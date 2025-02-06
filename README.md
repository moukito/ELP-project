# ELP-project

Ce projet contient trois sous-projets écrits dans différents langages de programmation :


1. [Projet ELM](#elm)
2. [projet Go](#go)
3. [Projet Javascript](#javascript)

---

## ELM

### Description

Le but de ce projet est de créer un langage TcTurtle qui permet à un utilisateur d'écrire un programme et de pourvoir ensuite observer le résultat directement.

### Installation
1. **Cloner le projet**
   ```shell
   git clone https://github.com/moukito/ELP-project.git
   cd ELP-project/elm
   ```
2. **Compiler le projet**
   ```shell
   elm make src/Main.elm --output Main.js
   ```
3. **Lancer le site web**
   - Sous Linux/MacOS :
   ```shell
   open index.html
   ```
   - Sous Windows :
   ```shell
   start index.html
   ```

---

## GO

### Description

Ce projet met en place un serveur tcp qui permet alors à différents clients d'envoyer une image qu'il souhaite numériser. L'image est alors traiter sur le serveur et renvoyé au client.

### Installation et Exécution
1. **Cloner le projet**
   ```shell
   git clone https://github.com/moukito/ELP-project.git
   cd ELP-project/go
   ```

2. **Compiler et exécuter**
   Compilez et exécutez l'application en utilisant les commandes Go :
   ```shell
   go build -o main main.go
   ./main
   ```

   Remarque : Assurez-vous d'avoir [Go](https://golang.org/) installé sur votre machine.

---

## Javascript

### Just One - Jeu en JavaScript

#### Description
Just One est une adaptation du jeu de société **Just One** en version numérique.  
C'est un jeu coopératif où les joueurs doivent faire deviner un mot mystère en proposant des indices uniques.

#### Règles du Jeu
1. Un joueur actif est désigné à chaque tour.
2. Il tire une carte et choisit un mot mystère parmi 5 propositions.
3. Les autres joueurs écrivent un indice secret pour l’aider à deviner le mot.
4. Les indices identiques sont annulés.
5. Le joueur actif tente de deviner le mot en utilisant les indices restants.
6. La partie continue jusqu'à l’épuisement de la pioche.

#### Installation
1. **Cloner le projet**
   ```shell
   git clone https://github.com/moukito/ELP-project.git
   cd ELP-project/js
   ```
2. **Installer les dépendances**
   ```shell
   npm install
   ```
3. **Lancer le serveur**
   ```shell
   npm start-server
   ```
4. **Lancer le client**
   ```shell
    npm start-client
    ```