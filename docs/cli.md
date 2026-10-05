"book" is the CLI for a book manager. It calls the REST APIs and
prints output in table format.

command format is like:

```
book <noun> <verb> [arguments] [flags]
```

Two nouns: book and collection.
book commands: add, list, show, edit, delete.
collection commands: add, list, show, edit, delete, add-book remove-book

book command examples:

1. book book add

format:
```
book book add --title <title> [--author <name>]... [--published <date>]
  [--edition <text>] [--genre <genre>] [--description <text>]
```

Only title is mandatory field and it supports multiple authors, example:

```
$ book book add --title "The Three Mistakes of My Life" \
    --author "Chetan Bhagat" \
    --author "Krish Malhotra" \
    --published 2008-01-01 \
    --genre "fiction"
ID           4
Title        The Three Mistakes of My Life
Authors      Chetan Bhagat, Krish Malhotra
Published    2008-01-01
Edition      -
Genre        fiction
Description  -
```

empty field prints as a hyphen. --published format is YYYY-MM-DD.


2. book book list

```
book book list [--author <name>] [--genre <genre>]
  [--published-after <date>] [--published-before <date>]
```

With no flags, lists every book:

```
$ book book list
ID  TITLE                              AUTHORS                        PUBLISHED   GENRE
4   The Three Mistakes of My Life      Chetan Bhagat, Krish Malhotra  2008-01-01  fiction
3   What Young India Wants             Chetan Bhagat                  2012-08-06  non fiction
2   A Game of Thrones                  George R. R. Martin            1996-08-01  fantasy
1   The Adventures of Sherlock Holmes  Arthur Conan Doyle             1890-01-01  fiction
```

Filters combine with AND. --author and --genre match one exact value; the date
filters bound the publication date and are inclusive on both ends:

```
$ book book list --genre fiction --published-after 2000-01-01 --published-before 2010-12-31
ID  TITLE                          AUTHORS                        PUBLISHED   GENRE
4   The Three Mistakes of My Life  Chetan Bhagat, Krish Malhotra  2008-01-01  fiction
```



3. book book show

```
book book show <id>
```

Shows one book with every field:

```
$ book book show 4
ID           4
Title        The Three Mistakes of My Life
Authors      Chetan Bhagat, Krish Malhotra
Published    2008-01-01
Edition      1st
Genre        fiction
Description  Three friends in Ahmedabad chase cricket, business and love.
```

4. book book edit

```
book book edit <id> --title <title> [--author <name>]... [--published <date>]
  [--edition <text>] [--genre <genre>] [--description <text>]
```

edit is a full replace: the book ends up with exactly the fields you pass, and any you omit are cleared.  --author replaces the whole author list, so to keep two authors and add a third, pass all three.

```
$ book book edit 4 --title "The Three Mistakes of My Life" \
    --author "Chetan Bhagat" --author "Krish Malhotra" \
    --published 2008-01-01 --edition "1st" --genre "fiction" \
    --description "Three friends in Ahmedabad chase cricket, business and love."
ID           4
Title        The Three Mistakes of My Life
Authors      Chetan Bhagat, Krish Malhotra
Published    2008-01-01
Edition      1st
Genre        fiction
Description  Three friends in Ahmedabad chase cricket, business and love.
```

5. book book delete

```
book book delete <id>
```

Removes the book and its links to any collections; the collections themselves are left
alone.

```
$ book book delete 4
Deleted book 4
```

collections exaples:

```
book collection show "To read"
```

1. book collection add

```
book collection add <name>
```

```
$ book collection add "To read"
Created collection "To read"
```

Fails if a collection with that name already exists.

2. book collection list

```
$ book collection list
NAME     BOOKS
Owned    0
To read  2
```

BOOKS is the no. of books in each collection.

3. book collection show

```
book collection show <name>
```


```
$ book collection show "To read"
To read  (2 books)

ID  TITLE                              AUTHORS              PUBLISHED   GENRE
2   A Game of Thrones                  George R. R. Martin  1996-08-01  fantasy
1   The Adventures of Sherlock Holmes  Arthur Conan Doyle   1890-01-01  fiction
```



4. book collection delete

```
book collection delete <name>
```

Removes the collection; the books in it are untouched.

```
$ book collection delete "To read"
Deleted collection "To read"
```

5. book collection add-book / remove-book

```
book collection add-book <name> <book-id>
book collection remove-book <name> <book-id>
```

```
$ book collection add-book "To read" 2
Added book 2 to "To read"

$ book collection remove-book "To read" 2
Removed book 2 from "To read"
```