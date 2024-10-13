const fs = require("node:fs")
const path = require("node:path");

const prismaFolderPath = path.resolve(__dirname, "../", process.argv[2]||"prisma")
const files = fs.readdirSync(prismaFolderPath).filter(s => s !== "schema.prisma" && s.endsWith(".prisma"))

let fileContent = "";

if (!files.includes("base.prisma")) throw "There is no \"base.prisma\""

fileContent += (fs.readFileSync(path.resolve(prismaFolderPath, "base.prisma"))+"").trim() + "\n\n"

for (let file of files) {
    if (file === "base.prisma") continue;
    fileContent += (fs.readFileSync(path.resolve(prismaFolderPath, file))+"").trim() + "\n\n"
}

fs.writeFileSync(path.resolve(prismaFolderPath, "schema.prisma"), fileContent.trim())