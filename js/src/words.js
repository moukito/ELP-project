const fs = require('fs');
const path = require('path');

/**
 * @file words.js
 * @description Gère le chargement des mots mystères à partir d'un fichier JSON.
 * Ce module fournit une méthode statique pour lire et récupérer la liste des mots
 * afin d'être utilisés dans le jeu Just One.
 */

class Words {
    /**
     * Charge la liste des mots depuis le fichier JSON.
     *
     * @async
     * @function loadWords
     * @returns {Promise<string[]>} Une promesse qui se résout avec un tableau contenant les mots mystères.
     * @throws {Error} Une erreur si le fichier JSON est introuvable ou mal formatté.
     *
     * @example
     * Words.loadWords()
     *   .then(words => console.log("Mots chargés:", words))
     *   .catch(error => console.error("Erreur de chargement:", error));
     */
    static loadWords() {
        return new Promise((resolve, reject) => {
            const filePath = path.join(__dirname, '../data/words.json');
            fs.readFile(filePath, 'utf8', (err, data) => {
                if (err) {
                    reject(new Error("Impossible de lire le fichier words.json. Assurez-vous qu'il existe."));
                } else {
                    try {
                        const json = JSON.parse(data);
                        if (!Array.isArray(json.words)) {
                            throw new Error("Le fichier words.json est mal formatté. Il doit contenir un tableau de mots.");
                        }
                        resolve(json.words);
                    } catch (e) {
                        reject(new Error("Erreur de parsing du fichier words.json : " + e.message));
                    }
                }
            });
        });
    }
}

module.exports = Words;
