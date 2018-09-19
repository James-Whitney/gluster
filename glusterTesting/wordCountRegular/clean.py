F = open("words.txt", "r")
book = F.read()
F.close()

book = book.replace('—', ' ')
chars = "*!(),.?:;“”’_"
for char in chars:
    book = book.replace(char, '')

result = open("wordsClean.txt", "w")
result.write(book.lower())
result.close()
print(book)