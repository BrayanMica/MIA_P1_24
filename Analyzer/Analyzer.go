package Analyzer

import (
	"MIA_P1_201907343/DiskManagement"
	"MIA_P1_201907343/FileSystem"
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

func Analyze() {

	for {
		var input string
		fmt.Print("-> ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()

		command, params := getCommandAndParams(input)

		fmt.Println("comando: ", command, "Parametros: ", params)

		AnalyzeCommnad(command, params)

		//mkdisk -size=3000 -unit=K -fit=BF
		//fdisk -size=300 -driveletter=A -name=Particion1
		//mount -driveletter=A -name=Particion1
		//mkfs -type=full -id=A119

	}
}

func AnalyzeCommnad(command string, params string) {

	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params)
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params)
	} else if strings.Contains(command, "mount") {
		fn_mount(params)
	} else if strings.Contains(command, "unmount") {
		fn_unmount(params)
	} else if strings.Contains(command, "mkfs") {
		fn_mkfs(params)
	} else {
		fmt.Println("Error: comando no encontrado")
	}

}

func fn_mkfs(input string) {
	// Define flags
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "", "Tipo")
	fs_ := fs.String("fs", "2fs", "Fs")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	FileSystem.Mkfs(*id, *type_, *fs_)

}

func fn_mount(input string) {
	// Define flags
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	driveletter := fs.String("driveletter", "", "Letra")
	name := fs.String("name", "", "Nombre")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	DiskManagement.Mount(*driveletter, *name)
}

func fn_fdisk(input string) {
	// Define flags
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	driveletter := fs.String("driveletter", "", "Letra")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "", "Unidad")
	type_ := fs.String("type", "", "Tipo")
	fit := fs.String("fit", "", "Ajuste")
	delete := fs.String("delete", "", "Eliminar")
	add := fs.String("add", "", "Agregar")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "driveletter", "name", "type":
			if flagValue == "" {
				// Aquí puedes poner el código que quieres ejecutar cuando flagValue es una cadena vacía
				fmt.Println("El parametro " + flagName + " no puede estar vacío")
				return
			} else {
				// Si flagValue no está vacío, entonces se establece el valor de la bandera
				fs.Set(flagName, flagValue)
			}
		case "delete", "add":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Parametro" + flagName + " no es valido")
		}
	}

	// Validate the flags
	if *size < 0 {
		fmt.Println("Error: Size debe de ser mayor que cero")
		return
	}

	// valitate driveletter
	*driveletter = strings.ToUpper(*driveletter)
	if !fileExists(*driveletter) {
		fmt.Println("Error: Drive letter does not exist")
		return
	}

	// validate name length
	if *name == "" {
		fmt.Println("Error: No has agregado un nombre en name")
		return
	}

	// validate unit
	*unit = strings.ToLower(*unit)
	if *unit == "" {
		*unit = "k"
	} else if *unit != "b" && *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit debe de ser 'b', 'k', or 'm'")
		return
	}

	// validate type
	*type_ = strings.ToUpper(*type_)
	if *type_ == "" {
		*type_ = "P"
	}
	if *type_ != "P" && *type_ != "E" && *type_ != "L" {
		fmt.Println("Error: Type debe de ser 'P', 'E', o 'L'")
		return
	}

	// validate fit
	*fit = strings.ToLower(*fit)
	if *fit == "" {
		*fit = "wf"
	} else if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: Fit debe de ser 'bf', 'ff', o 'wf'")
		return
	}

	// Call the function
	DiskManagement.Fdisk(*size, *driveletter, *name, *unit, *type_, *fit, *delete, *add)
}

// Verifica si un archivo existe
func fileExists(filepath string) bool {
	filepath = strings.ToUpper(filepath)
	filepath = filepath + ".bin"
	path := path.Join("./test/", filepath)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func fn_mkdisk(params string) {
	// Define flags
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "", "Unidad")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)

	// Track if we've seen the "size" flag
	sizeSeen := false

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit":
			if flagValue == "" {
				// Aquí puedes poner el código que quieres ejecutar cuando flagValue es una cadena vacía
				fmt.Println("El parametro " + flagName + " no puede estar vacío")
				return
			} else {
				// Si flagValue no está vacío, entonces se establece el valor de la bandera
				sizeSeen = true
				fs.Set(flagName, flagValue)
			}
		default:
			fmt.Println("Error: atributo no encontrado")
		}
	}

	// Check if we've seen the "size" flag
	if !sizeSeen {
		fmt.Println("Error: La bandera 'size' debe aparecer parametro obligatorio")
		return
	} else {
		// Call the function
		DiskManagement.Mkdisk(*size, *fit, *unit)
	}

}

func fn_rmdisk(params string) {
	// Define flags
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	driveletter := fs.String("driveletter", "", "Letradeunidad")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)

	errorRmdisk := false
	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: atributo no encontrado")
			errorRmdisk = true
			return
		}
	}

	if errorRmdisk {
		return
	} else {
		// Call the function
		DiskManagement.Rmdisk(*driveletter)
	}

}

func fn_unmount(params string) {
	println("unmount", params)
}
