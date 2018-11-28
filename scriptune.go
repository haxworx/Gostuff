package main

import "flag"
import "fmt"
import "os"
import "io/ioutil"
import "math/rand"
import "strings"
import "strconv"
import "time"
import "encoding/csv"

/* This is just a practice doing an old thing quickly */

const PROGRAM_NAME = "scriptune"

type Search struct {
	book_name string
	book_index int
	chapter_index int
	verse_index int
}

type Book struct {
	title    string
	chapters  []int
	index     int
}

type Bible struct {
	books map[string]Book
	filename string
	content []byte
	search Search
}

func helper() {
	fmt.Fprintf(os.Stderr, "Usage: %s <bookname> <chapter> <verse>\n", PROGRAM_NAME)
	fmt.Fprintf(os.Stderr, "OPTIONS:\n")
	fmt.Fprintf(os.Stderr, "    -r    random quote\n")

        os.Exit(0)
}

func _BookFromBlock(book *Book, block string) string {
	name_start := strings.Index(block, "title=")
        name_start += len("title=") + 1
	name_end := strings.Index(string(block[name_start:]), "\"")

	name_end += name_start
	name := block[name_start:name_end]

	chap_start := strings.Index(block, ">\n") + 2
	chap_end := strings.Index(string(block[chap_start:]),"\n")
	chap_end += chap_start

	verse_list := block[chap_start:chap_end]

	r := csv.NewReader(strings.NewReader(verse_list))

	book.title = name
	chapters, _ := r.ReadAll()
	for _, v := range chapters {
		for _, r := range v {
		i, _ := strconv.Atoi(r)
			book.chapters = append(book.chapters, i)
		}
	}

	return name
}

func ParseLayout(path string) map[string]Book {
	f, err := os.Open(path); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	books := make(map[string]Book)

	total, index := 0, 0

	bytes, _ := ioutil.ReadAll(f)
	contents := string(bytes)

	for {
		block_end := strings.Index(contents, "</book>")
                if block_end == -1 {
                        break
                } else {
			block := contents[0:block_end]

			book := Book {}
			book.index = index
			name := strings.ToUpper(_BookFromBlock(&book, block))
			books[name] = book;

			index++
			total += block_end + len("</book>")
			contents = string(bytes[total:])
		}
	}

	return books
}

func (b *Bible) Search(name string, chapter int, verse int) {
	s := b.books[strings.ToUpper(name)]

	if len(s.title) == 0 {
		fmt.Println("Book does not exist!")
		os.Exit(0)
	}

	if chapter > len(s.chapters) {
		fmt.Println("Book does not exist")
		os.Exit(0)
	}

	if verse > s.chapters[chapter - 1] {
		fmt.Println("Verse does not exist")
		os.Exit(0)
	}

	b.search.book_name = s.title
	b.search.book_index = s.index
	b.search.chapter_index = chapter
	b.search.verse_index = verse
}

func (b *Bible) ReadAll(filename string) []byte {
	f, err := os.Open(filename); if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(0)
	}

        b.content, _ = ioutil.ReadAll(f)
	f.Close()

	return b.content
}

func (b *Bible) Display() {
        bytes := b.content
	length := len(bytes)
	index := 0

	book, chapter, verse := -1, 0, 0

	for index < length {
		if bytes[index] == '@' {
			book++
			chapter = 0
			verse = 0
		} else if bytes[index] == '^' {
			chapter++
			verse = 0
		} else if bytes[index] == '~' {
			verse++
		} else if book == b.search.book_index && chapter == b.search.chapter_index && verse == b.search.verse_index {
			fmt.Printf("%c", bytes[index])
		}
		index++
	}

	fmt.Printf("\n\t- %s %d:%d\n", b.books[strings.ToUpper(b.search.book_name)].title, b.search.chapter_index, b.search.verse_index)
}

func (b *Bible) Random() {
	rand.Seed(int64(time.Now().UnixNano()))
        b.search.book_index = rand.Intn(len(b.books))

	book_name := ""

	for key := range b.books {
		if b.search.book_index == b.books[key].index {
			book_name = key
			break
		}
	}

	b.search.book_name = book_name

	rand.Seed(int64(time.Now().UnixNano()))
	b.search.chapter_index = rand.Intn(len(b.books[book_name].chapters))

	rand.Seed(int64(time.Now().UnixNano()))
	b.search.verse_index = rand.Intn(b.books[book_name].chapters[b.search.chapter_index])

	b.search.verse_index++
	b.search.chapter_index++
}

func BibleNew(layout_file, filename string) Bible {
	bible := Bible { }

        bible.books = ParseLayout(layout_file)
        bible.ReadAll(filename)

	return bible;
}

func main() {
	do_random := flag.Bool("r", false, "random verse")

	flag.Parse()

	if len(os.Args) != 4 && !*do_random {
		helper()
	}

	bible := BibleNew("kjv.layout", "kjv_nt.txt");

	switch (*do_random) {
	case true:
		bible.Random()
	case false:
		bookname := os.Args[1]
		chapter, _ := strconv.Atoi(os.Args[2])
		verse, _ := strconv.Atoi(os.Args[3])
		bible.Search(bookname, chapter, verse)
	}

	bible.Display();
}

