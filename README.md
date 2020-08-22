#NixPlayJpgTools
NAME: NixPlayJpgTools - Transforme un png en jpg ou lit des jpg pour les sauver non **entrelacés**. Sauve les fichiers dans le répertoire output avec sous directory. Si aucune commande est donnée alors les images sont testées sur leur taille et celles qui sont au dessus de 1280x720 (ou 720x1280 si portrait) sont listées.

USAGE: NixPlayJpgTools.exe [global options] command [command options] [arguments...]

VERSION: 1.0.0

COMMANDS: jpg relit des jpg et les sauve non entrelacés dans le répertoire output 
* **NixPlayJpgTools.exe jpg --resize** va convertir des images jpg en jpg non entrelacés et les retailler si besoin, le fichier de sortie sera en .jpg 
* **NixPlayJpgTools.exe --portrait jpg --resize** va faire comme ci-dessus mais va utiliser une taille 720x1280 png convertit un png en jpg et l'écrit dans le répertoire output 
* **NixPlayJpgTools.exe png --resize** va convertir des images png en jpg et les retailler si besoin 
* **NixPlayJpgTools.exe --portrait png --resize** va faire comme ci-dessus mais va utiliser une taille 720x1280 help, h Shows a list of commands or help for one command

Les jpg sont toujours non entrelacés.

GLOBAL OPTIONS: 
* --nowait Ne pas attendre une touche après process (default: false) 
* --portrait Va utiliser 702x1280 pour le resize si demandé (default: false) 
* --help, -h show help (default: false) 
* --version, -v print the version (default: false)

Licence MIT
