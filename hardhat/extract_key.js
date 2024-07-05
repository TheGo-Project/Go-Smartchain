const fs = require('fs');
const path = require('path');
const keythereum = require('keythereum');

// Получаем аргументы командной строки
const args = process.argv.slice(2);
const keystorePath = args[0];
const address = args[1];
const password = args[2];

if (!keystorePath || !address || !password) {
    console.error('Usage: node extract_private_key.js <keystorePath> <address> <password>');
    process.exit(1);
}

// Приводим адрес к нижнему регистру, так как файлы keystore могут использовать этот формат
const lowercasedAddress = address.toLowerCase();

// Найдем файл keystore, который содержит адрес
const files = fs.readdirSync(keystorePath);
const keyFile = files.find(file => file.includes(lowercasedAddress));

if (!keyFile) {
    console.error(`Keystore file not found for address: ${address}`);
    process.exit(1);
}

const keyObject = JSON.parse(fs.readFileSync(path.join(keystorePath, keyFile)));
const privateKey = keythereum.recover(password, keyObject);

console.log(`Private Key: 0x${privateKey.toString('hex')}`);