package Utilities

import (
	"MIA_P1_201907343/Structs"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Funtion to create bin file
func CreateFile(name string) error {
	//Ensure the directory exists
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Create file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funtion to open bin file in read/write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

// Function to Write an object in a bin file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Function to Read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

// Function to read EBR from file
func ReadEBR(file *os.File, position int64) (Structs.EBR, error) {
	var ebr Structs.EBR
	// Move the cursor to the position
	_, err := file.Seek(position, 0)
	if err != nil {
		return ebr, err
	}
	// Read the EBR
	ebrData := make([]byte, binary.Size(ebr))
	_, err = file.Read(ebrData)
	if err != nil {
		return ebr, err
	}
	// Convert the bytes to EBR
	buffer := bytes.NewBuffer(ebrData)
	err = binary.Read(buffer, binary.LittleEndian, &ebr)
	if err != nil {
		return ebr, err
	}
	return ebr, nil
}

func DeletePartition(file *os.File, name string) error {
	// Load the MBR from the file
	var mbr Structs.MRB
	err := ReadObject(file, &mbr, 0)
	if err != nil {
		fmt.Println("Err DeletePartition==", err)
		return err
	}

	// Create a new array
	var newPartitions [4]Structs.Partition

	// Find the partition in the MBR
	for i, partition := range mbr.Partitions {
		if strings.TrimRight(string(partition.Name[:]), "\x00") == name {
			fmt.Println("Eliminando la partición:", name)

			// Copy the partitions before the one to delete
			copy(newPartitions[:i], mbr.Partitions[:i])

			// Copy the partitions after the one to delete
			copy(newPartitions[i:], mbr.Partitions[i+1:])

			// Replace the old partitions with the new ones
			mbr.Partitions = newPartitions

			// Write the updated MBR back to the file
			err = WriteObject(file, &mbr, 0)
			if err != nil {
				fmt.Println("Err DeletePartition==", err)
				return err
			}

			break
		}
	}

	fmt.Println("Particiones restantes:")
	for _, partition := range mbr.Partitions {
		fmt.Println(strings.TrimRight(string(partition.Name[:]), "\x00"))
	}

	return nil
}

func ModifyPartition(file *os.File, name string, change int, unit string) error {

	// Load the MRB from the file
	var mrb Structs.MRB
	err := ReadObject(file, &mrb, 0)
	if err != nil {
		return err
	}

	// validate unit
	unit = strings.ToUpper(unit)
	if unit == "" {
		return errors.New("unit no puede estar vacío")
	} else if unit == "B" {
		// El valor de change no cambia por que esta ingresado en bytes
	} else if unit == "K" {
		change = change * 1024
	} else if unit == "M" {
		change = change * 1024 * 1024
	} else if unit != "B" && unit != "K" && unit != "M" {
		return errors.New("unit debe ser 'B', 'K', o 'M'")
	}

	// calculo de espacio en memoria disponible de la particion
	var freeSpace int32

	for _, partition := range mrb.Partitions {
		freeSpace += partition.Size
	}
	freeSpace = mrb.MbrSize - freeSpace
	fmt.Println("Espacio libre: ", freeSpace/1024/1024, "MB")

	// Find the partition in the MRB
	for i, partition := range mrb.Partitions {
		if strings.TrimRight(string(partition.Name[:]), "\x00") == name {
			// Check if there is enough space to add or remove

			if change > 0 && freeSpace < int32(change) {
				return errors.New("no hay suficiente espacio libre para agregar a la partición")
			} else if change <= 0 && partition.Size <= -int32(change) {
				return errors.New("no hay suficiente espacio en la partición para quitar")
			}

			// Add or remove the space
			partition.Size += int32(change)
			freeSpace -= int32(change)
			// Update the partition in the MRB
			mrb.Partitions[i] = partition

			// Write the updated MRB back to the file
			err = WriteObject(file, &mrb, 0)
			if err != nil {
				return err
			}

			break
		}
	}

	// mostrar el nuevo valor size de cada particion

	fmt.Println("Tamaño Actualizado de particiones en MB:")
	for _, partition := range mrb.Partitions {
		fmt.Println(strings.TrimRight(string(partition.Name[:]), "\x00"), partition.Size/1024/1024, "MB")
	}
	fmt.Println("Nuevo Espacio libre: ", freeSpace/1024/1024, "MB")

	return nil
}
