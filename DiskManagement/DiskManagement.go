package DiskManagement

// importacion de librerias
import (
	"MIA_P1_201907343/Structs"
	"MIA_P1_201907343/Utilities"
	"bytes"
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
	var studentID = "201907343"
	// convertir driveletter a mayusculas
	driveletter = strings.ToUpper(driveletter)
	fmt.Println("======Inicio MOUNT======")
	fmt.Println("Driveletter:", driveletter)
	fmt.Println("Name:", name)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: el archivo del disco no existe")
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: no se pudo leer el MBR del disco")
		return
	}

	// Find the partition with the given name
	var partition *Structs.Partition
	var partitionIndex int
	for i := 0; i < 4; i++ {
		partitionName := strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")
		if partitionName == name {
			partition = &TempMBR.Partitions[i]
			partitionIndex = i
			break
		}
	}

	if partition == nil {
		fmt.Println("Error: no se encontró la partición con el nombre dado")
		return
	}

	// Mount the partition
	partition.Status[0] = '1' // Change status to mounted

	// Update ID
	id := []byte(driveletter + strconv.Itoa(partitionIndex+1) + studentID[len(studentID)-2:])
	copy(partition.Id[:], id[:4])

	// Write the updated MBR back to the file
	file.Seek(0, 0)
	binary.Write(file, binary.LittleEndian, &TempMBR)

	fmt.Println("Partición montada correctamente con ID:", string(partition.Id[:]))
	fmt.Println("======Fin MOUNT======")
	fmt.Println("")
}

func Unmount(partitionID string) {
	fmt.Println("======Inicio UNMOUNT======")
	fmt.Println("Partition ID:", partitionID)

	// Open bin file
	driveletter := string(partitionID[0])
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: el archivo del disco no existe")
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: no se pudo leer el MBR del disco")
		return
	}

	// Find the partition with the given ID
	var partition *Structs.Partition
	for i := 0; i < 4; i++ {
		targetPartitionIDpartitionID := strings.TrimRight(string(TempMBR.Partitions[i].Id[:]), "\x00")
		if partitionID == targetPartitionIDpartitionID {
			partition = &TempMBR.Partitions[i]
			break
		}
	}

	if partition == nil {
		fmt.Println("Error: no se encontró la partición con el ID dado")
		return
	}

	// Unmount the partition
	partition.Status[0] = '0' // Change status to unmounted

	// Write the updated MBR back to the file
	file.Seek(0, 0)
	binary.Write(file, binary.LittleEndian, &TempMBR)

	fmt.Println("Partición desmontada correctamente con ID:", partitionID)
	fmt.Println("======Fin UNMOUNT======")
	fmt.Println("")
}

func Fdisk(size int, driveletter string, name string, unit string, type_ string, fit string) {
	fmt.Println("======Inicio FDISK======")
	//fmt.Println("Size:", size)
	//fmt.Println("Driveletter:", driveletter)
	//fmt.Println("Name:", name)
	//fmt.Println("Unit:", unit)
	//fmt.Println("Type:", type_)
	//fmt.Println("Fit:", fit)

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

	// Calculate the total space occupied by existing partitions
	var totalPartitionSize int32 = 0
	for _, partition := range TempMBR.Partitions {
		totalPartitionSize += partition.Size
	}

	// Check if there is enough space left on the disk
	if TempMBR.MbrSize-totalPartitionSize < int32(size) {
		fmt.Println("Espacio total en disco:", TempMBR.MbrSize)
		fmt.Println("Espacio ocupado por particiones:", totalPartitionSize)
		fmt.Println("Espacion disponible:", TempMBR.MbrSize-totalPartitionSize)
		fmt.Println("Espacio requerido:", size)
		fmt.Println("Error: No hay suficiente espacio en el disco")
		fmt.Println("======Fin FDISK======")
		fmt.Println("")
		return
	}

	// check if all partitions are full

	allFull := true

	// Count primary and extended partitions
	var primaryCount, extendedCount int = 0, 0
	for _, partition := range TempMBR.Partitions {
		if partition.Size == 0 {
			allFull = false
			break
		} else if strings.TrimRight(string(partition.Name[:]), "\x00") == name {
			fmt.Println("Error: Ya existe una partición con el nombre: ", name)
			fmt.Println("======Fin FDISK======")
			fmt.Println("")
			return
		} else if partition.Type[0] == byte('P') {
			primaryCount++
		} else if partition.Type[0] == byte('E') {
			//fmt.Println("Partition Encontrada: ", string(partition.Name[:]), " ", string(partition.Type[:]))
			extendedCount++
		}
	}

	if allFull && (type_ == "P" || type_ == "E") {
		fmt.Println("Error: Todas los espacios de particiones están llenas")
		fmt.Println("======Fin FDISK======")
		fmt.Println("")
		return
	}

	// Check if there is more than one extended partition
	if extendedCount == 1 && type_ == "E" {
		fmt.Println("Error: Solo puede haber una partición extendida por disco.")
		fmt.Println("======Fin FDISK======")
		fmt.Println("")
		return
	}

	// check if all partitions are not primary
	if primaryCount == 3 && type_ == "P" {
		fmt.Println("Error: No puedes agregar todas las particiones del tipo primario debes de agregar por lo menos 1 del tipo Extendida.")
		fmt.Println("======Fin FDISK======")
		fmt.Println("")
		return
	}

	// cheacl if patition == "L"
	if extendedCount == 0 && type_ == "L" {
		fmt.Println("Error: No se puede agregar una particion de tipo logica sin tener una extendida")
		fmt.Println("======Fin FDISK======")
		fmt.Println("")
		return
	}

	if type_ == "E" {
		// Encontrar un espacio libre en el MBR
		var freeIndex int
		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size == 0 {
				freeIndex = i
				break
			}
		}

		// Configurar la partición extendida en el MBR
		TempMBR.Partitions[freeIndex].Size = int32(size)
		TempMBR.Partitions[freeIndex].Start = int32(binary.Size(TempMBR.Partitions)) // Empieza después de las particiones en el MBR
		copy(TempMBR.Partitions[freeIndex].Name[:], name)
		copy(TempMBR.Partitions[freeIndex].Fit[:], fit)
		copy(TempMBR.Partitions[freeIndex].Status[:], "0")
		copy(TempMBR.Partitions[freeIndex].Type[:], type_)
		TempMBR.Partitions[freeIndex].Correlative = 1 // Para la primera partición extendida

		// Escribir el MBR actualizado en el archivo
		if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
			return
		}

		// Crear el EBR virtual
		var TempEBR Structs.EBR
		TempEBR.PartFit = fit[0]
		TempEBR.PartSize = int32(size)
		TempEBR.PartNext = -1 // No hay siguiente EBR
		copy(TempEBR.PartName[:], name)

		// Calcular el inicio de la partición
		var start int32 = 0
		for i := 0; i < freeIndex; i++ {
			start += TempMBR.Partitions[i].Size
		}
		TempEBR.PartStart = start

		// Convertir la estructura EBR en una secuencia de bytes
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, &TempEBR)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
			return
		}

		// Escribir la secuencia de bytes en el archivo
		_, err = file.WriteAt(buf.Bytes(), int64(TempEBR.PartStart))
		if err != nil {
			fmt.Println("WriteAt failed:", err)
			return
		}

		fmt.Println("Partición extendida creada correctamente")
		fmt.Println("======Fin FDISK======")
		fmt.Println("")
		return
	}

	if type_ == "L" {
		// Encontrar la partición extendida
		var extendedPartition *Structs.Partition
		var i int
		for i = 0; i < 4; i++ {
			if string(TempMBR.Partitions[i].Type[:]) == "E" {
				extendedPartition = &TempMBR.Partitions[i]
				break
			}
		}

		if extendedPartition == nil {
			fmt.Println("No se encontró una partición extendida para agregar la partición lógica")
			return
		}

		// Crear un nuevo EBR
		var newEBR Structs.EBR
		newEBR.PartFit = 'F'          // Ajuste por defecto
		newEBR.PartSize = int32(size) // Tamaño de la partición
		newEBR.PartStart = extendedPartition.Start + extendedPartition.Size
		newEBR.PartMount = '0' // Estado por defecto
		newEBR.PartNext = -1   // No hay siguiente EBR

		// Escribir el nuevo EBR en el archivo
		file.Seek(int64(newEBR.PartStart), 0)
		binary.Write(file, binary.LittleEndian, &newEBR)

		// Actualizar el MBR
		TempMBR.Partitions[i].Start = newEBR.PartStart
		TempMBR.Partitions[i].Size = newEBR.PartSize
		file.Seek(0, 0)
		binary.Write(file, binary.LittleEndian, &TempMBR)

		fmt.Println("Partición lógica creada correctamente")
		fmt.Println("======Fin FDISK======")
		fmt.Println("")
		return
	}
	// Print object
	//Structs.PrintMBR(TempMBR)

	if type_ == "P" {

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
		//Structs.PrintMBR(TempMBR2)

		// Close bin file
		err = file.Close()
		if err != nil {
			//manejar el error
			fmt.Println("Error: ", err)
			return
		}

		fmt.Println("Partición primaria creada correctamente")
		fmt.Println("======Fin FDISK======")
		fmt.Println("")
	}
}

func Fdisk_Delete(delete string, name string, driveletter string) {

	fmt.Println("======Inicio FDISK_DELETE======")
	fmt.Println("Delete:", delete)
	fmt.Println("Name:", name)
	fmt.Println("Driveletter:", driveletter)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	if err := Utilities.DeletePartition(file, name); err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("======Fin FDISK_DELETE======")
		fmt.Println("")
		return
	}
	// fdisk -delete=full -name=Part2 -driveletter=A

	fmt.Println("======Fin FDISK_DELETE======")
	fmt.Println("")
}

func Fdisk_Add(add int, unit string, name string, driveletter string) {
	fmt.Println("======Inicio FDISK_ADD======")
	fmt.Println("Add:", add)
	fmt.Println("Unit:", unit)
	fmt.Println("Name:", name)
	fmt.Println("Driveletter:", driveletter)

	// Open bin file
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	// fdisk -add=1000 -unit=k -driveletter=A -name=Part26

	// verificacion si la particion que estoy agregando o quitando espacio existe en el disco
	var TempMBR Structs.MRB

	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	partitionExists := false
	for i := 0; i < 4; i++ {
		if strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00") == name {
			partitionExists = true
		}
	}

	if partitionExists {
		// LLama a la funcion ModifyPartition para agregar espacio a la particion
		err = Utilities.ModifyPartition(file, name, add, unit)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	} else {
		fmt.Println("Error: La particion no existe")
	}
	// Close bin file
	err = file.Close()
	if err != nil {
		fmt.Println("Error cerrando el archivo: ", err)
		return
	}
	fmt.Println("======Fin FDISK_ADD======")
	fmt.Println("")
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
		//
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

	// Asignando buffer de 1048 bytes
	buffer := make([]byte, 1024)
	for i := range buffer {
		buffer[i] = 0
	}

	// create array of byte(0)
	for i := 0; i < size; i += len(buffer) {
		err := Utilities.WriteObject(file, buffer, int64(i))
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
	//Structs.PrintMBR(TempMBR)

	// Close bin file
	defer file.Close()

	fmt.Println("======Fin MKDISK======")
	fmt.Println("")
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
	fmt.Println("")
}

func Execute(path string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("El programa se está ejecutando en el directorio: ", dir)

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

func Rep(id string, path string) {

	fmt.Println("======Inicio REP======")
	print(path)
	path = "/home/estiben/Documentos/go/src/Proyectos/MIA_P1_201907343/Reportes base"
	// imprimir el nombre de cada disco de la carpeta test
	files, err := ioutil.ReadDir("./test")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for _, file := range files {
		fmt.Println(file.Name())
		// leer mbr de cada disco
		filepath := "./test/" + file.Name()
		file, err := Utilities.OpenFile(filepath)
		if err != nil {
			return
		}

		var TempMBR Structs.MRB
		// Read object from bin file
		if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
			return
		}

		// enviar mrbs a la funcion GenerateDotFile
		if err := GenerateDotFile(&TempMBR, path); err != nil {
			fmt.Println("Error: ", err)
		}

	}
	fmt.Println("======Fin REP======")
}
func GenerateDotFile(mbr *Structs.MRB, filepath string) error {
	// Inicia el código DOT
	dotCode := "digraph G {\n"
	dotCode += "node [shape=record];\n"
	dotCode += "MBR [label=\"{MBR"

	// Añade cada partición al código DOT
	totalSize := mbr.MbrSize
	for i, partition := range mbr.Partitions {
		partitionID := strings.TrimRight(string(partition.Id[:]), "\x00")
		partitionSize := partition.Size
		percentage := float64(partitionSize) / float64(totalSize) * 100
		dotCode += fmt.Sprintf("|<f%d> %s (%.2f%%)", i, partitionID, percentage)
	}

	dotCode += "}\"];\n"
	dotCode += "}\n"

	// Escribe el código DOT a un archivo
	return ioutil.WriteFile(filepath, []byte(dotCode), 0644)
}
