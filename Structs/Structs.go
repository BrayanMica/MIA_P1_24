package Structs

import (
	"fmt"
)

//  =============================================================

type MRB struct {
	MbrSize      int32
	CreationDate [10]byte
	Signature    int32
	Fit          [1]byte
	Partitions   [4]Partition
}

func PrintMBR(data MRB) {
	fmt.Printf("CreationDate: %s, fit: %s, size: %d\n", string(data.CreationDate[:]), string(data.Fit[:]), data.MbrSize)
	for i := 0; i < 4; i++ {
		PrintPartition(data.Partitions[i])
	}
}

//  =============================================================

type Partition struct {
	Status      [1]byte  // indica si la particion esta montada o no
	Type        [1]byte  // indica el tipo de particion, primaria(P) o extendida(E)
	Fit         [1]byte  // Tipo de ajuste de la particion (Best(B), First(F), Worst(W))
	Start       int32    // Indica en que byte del disco inicia la particion
	Size        int32    // Tamaño de la particion en bytes
	Name        [16]byte // Nombre de la particion
	Correlative int32    // Correlativo de la particion es un numero
	Id          [4]byte  // Identificador unico de la particion
}

func PrintPartition(data Partition) {
	fmt.Printf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s\n", string(data.Name[:]), string(data.Type[:]), data.Start, data.Size, string(data.Status[:]), string(data.Id[:]))
}

//  =============================================================

type EBR struct {
	PartMount byte     // Indica si la partición está montada o no
	PartFit   byte     // Tipo de ajuste de la partición. Tendrá los valores B (Best), F (First) o W (worst)
	PartStart int32    // Indica en qué byte del disco inicia la partición
	PartSize  int32    // Contiene el tamaño total de la partición en bytes.
	PartNext  int32    // Byte en el que está el próximo EBR. -1 si no hay siguiente
	PartName  [16]byte // Nombre de la partición
}

//  =============================================================

type Superblock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [17]byte
	S_umtime            [17]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_fist_ino          int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

//  =============================================================

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [17]byte
	I_ctime [17]byte
	I_mtime [17]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

//  =============================================================

type Fileblock struct {
	B_content [64]byte
}

//  =============================================================

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type Folderblock struct {
	B_content [4]Content
}

//  =============================================================

type Pointerblock struct {
	B_pointers [16]int32
}

//  =============================================================

type Content_J struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [17]byte
}

type Journaling struct {
	Size      int32
	Ultimo    int32
	Contenido [50]Content_J
}
