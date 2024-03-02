# Proyecto de Manejo e Implementación de Archivos

## Universidad de San Carlos de Guatemala
### Facultad de Ingeniería
#### Escuela de Ciencias y Sistemas

### Introducción
El proyecto de Manejo e Implementación de Archivos es parte del curso impartido en la Universidad de San Carlos de Guatemala, diseñado para que los estudiantes adquieran conocimientos en la administración de archivos, tanto a nivel de hardware como de software, explorando sistemas de archivos, particiones y otros conceptos esenciales. El objetivo principal es que los estudiantes apliquen estos conocimientos en el desarrollo de un proyecto que les permita comprender cada uno de los temas impartidos en la clase magistral y en el laboratorio. Además, este proyecto sienta las bases para futuros cursos relacionados, como las bases de datos.

### Objetivos
- Aprender a administrar archivos y escribir estructuras en Go.
- Comprender el sistema de archivos EXT3 y EXT2.
- Aplicar el formateo rápido y completo en una partición.
- Crear una aplicación de comandos.
- Aplicar la teoría de ajustes y particiones.
- Utilizar Graphviz para mostrar reportes.
- Restringir y administrar el acceso a los archivos y carpetas en ext3/ext2 por medio de usuarios.
- Administrar los usuarios y permisos por medio de grupo.

### Discos
Los discos serán simulados mediante archivos binarios con extensión `.dsk`, donde se almacenarán las diferentes estructuras para simular el funcionamiento del almacenamiento de un sistema de archivos y las particiones que pueda contener.

### Master Boot Record (MBR)
El MBR, ubicado en el primer sector del disco, proporciona información del sistema de archivos y de las particiones. Contendrá los siguientes valores:

- `mbr_tamano`: Tamaño total del disco en bytes.
- `mbr_fecha_creacion`: Fecha y hora de creación del disco.
- `mbr_dsk_signature`: Número aleatorio que identifica de forma única a cada disco.
- `dsk_fit`: Tipo de ajuste de la partición.
- `mbr_partitions`: Estructura con información de las 4 particiones.

### Partition
Una partición es una división lógica de un disco tratada como una unidad separada. Contendrá información sobre su estado, tipo, ajuste, posición inicial, tamaño, nombre y correlativo.

### Extended Boot Record (EBR)
El EBR actúa como descriptor de una unidad lógica y apunta al espacio donde se escribirá el siguiente EBR. Se considera una especie de lista enlazada que contiene información sobre la partición lógica.

### Sistema de Archivos Ext2/Ext3
Se deberán implementar las estructuras de los sistemas de archivos EXT2 y EXT3, incluyendo el súper bloque, inodos y bloques, así como el formato de las particiones.

### Estructuras para Carpetas y Archivos
- **Súper Bloque:** Contiene información sobre el sistema de archivos, como el tipo de sistema, el número total de inodos y bloques, la cantidad de bloques e inodos libres, entre otros.
- **Inodos (index node):** Contiene información sobre un archivo o carpeta, incluyendo el tamaño, permisos, fechas de acceso, creación y modificación, y los bloques asignados.

Este proyecto proporciona una oportunidad para que los estudiantes exploren y apliquen conceptos fundamentales relacionados con el manejo de archivos y sistemas de archivos, lo que sentará las bases para su desarrollo profesional en el campo de la informática y la ingeniería de sistemas.
