package main

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// 1.利用fmt.Fscan和fmt.Fprint系列函数进行读写
type Player struct {
	name   string
	age    int
	gender string
}

func (p Player) String() string {
	return fmt.Sprintf("%s %d %s", p.name, p.age, p.gender)
}

func directWriteADTToFile(path string, players []Player) error {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("open file error:", err)
		return err
	}
	defer func() {
		f.Sync()
		f.Close()
	}()

	for _, player := range players {
		_, err = fmt.Fprintf(f, "%s\n", player)
		if err != nil {
			return err
		}
	}
	return nil
}

// 2.利用binary.Read和binary.Write函数进行读写
type PlayerBinary struct {
	Name   [20]byte
	Age    int16
	Gender [6]byte
}

func directBinaryWriteADTToFile(path string, players []PlayerBinary) error {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("open file error:", err)
		return err
	}
	defer func() {
		f.Sync()
		f.Close()
	}()

	for _, player := range players {
		err = binary.Write(f, binary.BigEndian, &player)
		if err != nil {
			return err
		}
	}
	return nil
}

// 利用gob包读取写图
func directGobWriteADTToFile(path string, players []Player) error {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("open file error:", err)
		return err
	}
	defer func() {
		f.Sync()
		f.Close()
	}()
	enc := gob.NewDecoder(f)

	for _, player := range players {
		err = enc.Decode(player)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var players = []Player{
		{"Tommy", 18, "male"},
		{"Lucy", 17, "female"},
		{"George", 19, "male"},
	}
	err := directWriteADTToFile("players.dat", players)
	if err != nil {
		fmt.Println("write file error: ", err)
		return
	}
	f, err := os.Open("players.dat")
	if err != nil {
		fmt.Println("open file error: ", err)
		return
	}

	var player Player
	for {
		_, err := fmt.Fscanf(f, "%s %d %s", &player.name, &player.age, &player.gender)
		if err == io.EOF {
			fmt.Println("read meet EOF")
			return
		}
		if err != nil {
			fmt.Println("read file error: ", err)
			return
		}
		fmt.Printf("%s %d %s\n", player.name, player.age, player.gender)
	}
	// 我们日常用得最多的格式化输出/输入函数Printf系列和Scanf
	// 系列本质上是当io.Writer为os.Stdout或io.Reader为os.Stdin时的特例

	var playerBinarys [3]PlayerBinary

	copy(playerBinarys[0].Name[:], []byte("Tommy"))
	playerBinarys[0].Age = 18
	copy(playerBinarys[0].Gender[:], []byte("male"))

	playerBinarys[1].Age = 20
	copy(playerBinarys[1].Name[:], []byte("Jack"))
	copy(playerBinarys[1].Gender[:], []byte("male"))

	playerBinarys[2].Age = 17
	copy(playerBinarys[2].Name[:], []byte("Lily"))
	copy(playerBinarys[2].Gender[:], []byte("female"))

	errBinary := directBinaryWriteADTToFile("playerBinarys.dat", playerBinarys[:])
	if errBinary != nil {
		fmt.Println("write file error: ", err)
		return
	}

	fBinary, errBinary := os.Open("playerBinarys.dat")
	if errBinary != nil {
		fmt.Println("open file error: ", err)
		return
	}

	var playerBinary PlayerBinary
	for {
		binary.Read(fBinary, binary.BigEndian, &playerBinary)
		if err == io.EOF {
			fmt.Println("read meet EOF")
			return
		}
		if err != nil {
			fmt.Println("read file error: ", err)
			return
		}
		fmt.Printf("%s %d %s\n", playerBinary.Name, playerBinary.Age, playerBinary.Gender)
	}
}
