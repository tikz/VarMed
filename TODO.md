## Features
- abSwitch aparentemente escribe a la base de datos que trae, lo que hace que no se pueda paralelizar porque se pisan.
  El mutex genera cuello de botella, mas que es un paso lento. Considerar hacer una copia de abswitch en una carpeta aparte para cada trabajo.
  Ya probe. Explota la RAM.

## Features menores QoL
- Flag de guardar o no archivos resultantes en config.yaml, ahora estoy comentando y recompilando dependiendo si es local o prod
- Archivo log de solo errores

## Refactor
- Estructurar en subpackages de forma diferente
- Ahora que hay estructuras de datos mas o menos no tan cambiantes, redefinir la forma de manejar channels en pipeline.go
- SIFTS mapping: es el JSON tal cual como viene y se usa en varios lados, la estructura es incomoda.
- Feo feo el regexp hell. En UniProt TXT creo que no se justifica.
- Unificar funciones que lanzan procesos y leen archivos, veo logica repetida.
- Ser explicito en que tipo de posiciones devuelve la API JSON en cada caso, si uniprot o pdb. Ahora hay una mezcla que yo solo la sé.