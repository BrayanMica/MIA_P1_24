package DiskManagement

// importacion de librerias
import (
	"MIA_P1_201907343/Structs"
	"MIA_P1_201907343/Utilities"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// variables globales
var fileName string = ""

func Mount(driveletter string, name string) {
	fmt.Println("======Inicio MOUNT======")
	fmt.Println("Driveletter:", driveletter)
	fmt.Println("Name:", name)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)

	fmt.Println("-------------")

	var index int = -1
	var count = 0
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			count++
			if strings.Contains(string(TempMBR.Partitions[i].Name[:]), name) {
				index = i
				break
			}
		}
	}

	if index != -1 {
		fmt.Println("Partition Encontrada:")
		Structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition no Encontrada")
		return
	}

	// id = DriveLetter + Correlative + 19

	id := strings.ToUpper(driveletter) + strconv.Itoa(count) + "19"

	copy(TempMBR.Partitions[index].Status[:], "1")
	copy(TempMBR.Partitions[index].Id[:], id)

	// Overwrite the MBR
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR2)

	// Close bin file
	defer file.Close()

	fmt.Println("======Fin MOUNT======")
}

func Fdisk(size int, driveletter string, name string, unit string, type_ string, fit string, delete string, add string) {
	fmt.Println("======Inicio FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Driveletter:", driveletter)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)

	fmt.Println("-------------")

	var count = 0
	var gap = int32(0)
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			count++
			gap = TempMBR.Partitions[i].Start + TempMBR.Partitions[i].Size
		}
	}

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size == 0 {
			TempMBR.Partitions[i].Size = int32(size)

			if count == 0 {
				TempMBR.Partitions[i].Start = int32(binary.Size(TempMBR))
			} else {
				TempMBR.Partitions[i].Start = gap
			}

			copy(TempMBR.Partitions[i].Name[:], name)
			copy(TempMBR.Partitions[i].Fit[:], fit)
			copy(TempMBR.Partitions[i].Status[:], "0")
			copy(TempMBR.Partitions[i].Type[:], type_)
			TempMBR.Partitions[i].Correlative = int32(count + 1)
			break
		}
	}

	// Overwrite the MBR
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR2)

	// Close bin file
	err = file.Close()
	if err != nil {
		//manejar el error
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("======Fin FDISK======")
}

func Mkdisk(size int, fit string, unit string) {
	fmt.Println("======Inicio MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	// validate fit equals to b/w/f
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit debe de ser bf, wf o ff")
		return
	}

	// validate size > 0
	if size <= 0 {
		fmt.Println("Error: Size debe de ser mayor que 0")
		return
	}

	// validate unit equals to k/m
	if unit == "" {
		unit = "m"
	}
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Unit debe de ser k o m")
		return
	}

	createNextFile("./test/")
	// Create file
	// err := Utilities.CreateFile("./test/" + filename)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }

	// Set the size in bytes
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// Open bin file
	file, err := Utilities.OpenFile("./test/" + fileName)
	if err != nil {
		return
	}

	// Write 0 binary data to the file

	// create array of byte(0)
	for i := 0; i < size; i++ {
		err := Utilities.WriteObject(file, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// Create a new instance of MRB
	var newMRB Structs.MRB
	newMRB.MbrSize = int32(size)
	newMRB.Signature = 10 // random
	copy(newMRB.Fit[:], fit)
	currentDate := time.Now().Format("2006-01-02")
	copy(newMRB.CreationDate[:], currentDate)

	// Write object in bin file
	if err := Utilities.WriteObject(file, newMRB, 0); err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)

	// Close bin file
	defer file.Close()

	fmt.Println("======Fin MKDISK======")

}

func createNextFile(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fileNames := make(map[string]bool)
	for _, file := range files {
		name := strings.TrimSuffix(file.Name(), ".bin")
		fileNames[name] = true
	}

	nextFileName := "A"
	for fileNames[nextFileName] {
		nextFileName = incrementFileName(nextFileName)
	}
	nextFileName += ".bin"
	fileName = nextFileName
	_, err = os.Create(path + nextFileName)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func incrementFileName(fileName string) string {
	if fileName == "" {
		return "A"
	}
	lastChar := fileName[len(fileName)-1]
	if lastChar < 'Z' {
		return fileName[:len(fileName)-1] + string(lastChar+1)
	}
	return incrementFileName(fileName[:len(fileName)-1]) + "A"
}

func Rmdisk(driveletter string) {
	fmt.Println("======Inicio RMDISK======")
	fmt.Println("Driveletter:", driveletter)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	err := os.Remove(filepath)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("======Fin RMDISK======")
}

func Execute(path string) {
	fmt.Println("======Inicio FILESYSTEM======")
	fmt.Println("Path:", path)

	// Open bin file
	file, err := Utilities.OpenFile(path)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Structs.PrintMBR(TempMBR)

	// Close bin file
	defer file.Close()

	fmt.Println("======Fin FILESYSTEM======")
}
