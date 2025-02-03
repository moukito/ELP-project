const fs = require('fs');
const path = require('path');

class Words {
    /**
     * Charge la liste des mots depuis le fichier JSON.
     * @returns {Promise<Array>} Une promesse qui se résout avec un tableau de mots.
     */
    static loadWords() {
        return new Promise((resolve, reject) => {
            const filePath = path.join(__dirname, '../data/words.json');
            fs.readFile(filePath, 'utf8', (err, data) => {
                if (err) {
                    reject(err);
                } else {
                    try {
                        const json = JSON.parse(data);
                        resolve(json.words);
                    } catch (e) {
                        reject(e);
                    }
                }
            });
        });
    }
}

module.exports = Words;
