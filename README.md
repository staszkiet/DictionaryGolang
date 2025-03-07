# Dictionary App

## Requirements
- Docker with docker-compose
- Go language (if you want to use the client app)

## Start the server and DB

1. Clone the repository
2. Go into the server folder
3. Run `docker-compose build`
4. Run `docker-compose up -d`

## Start the client (optional)

1. Go into the app folder
2. Run `go mod tidy`
3. Run `go run .`

## Queries and mutations examples

### Create polish-english translation

**GraphQL:**
```graphql
mutation createWord {
  createWord(
    polish: "rower"
    translation: {
      english: "bike"
      sentences: [
        "My bike is green."
        "I like my bike."
      ]
    }
  )
}
```

**Client:**
```
ADD rower bike (My bike is green.) (I like my bike.)
```

### Add a translation to an existing polish word

**GraphQL:**
```graphql
mutation {
  createTranslation(
    polish: "rower"
    translation: {
      english: "bicycle"
      sentences: [
        "I like my bicycle."
        "The bicycle is under repair."
      ]
    }
  )
}
```

**Client:**
```
ADD_TRANSLATION rower bicycle (I like my bicycle) (The bicycle is under repair)
```

### Add example sentence to an existing translation

**GraphQL:**
```graphql
mutation {
  createSentence(
    polish: "rower"
    english: "bicycle"
    sentence: "I dont like my bicycle."
  )
}
```

**Client:**
```
ADD_SENTENCE rower bicycle (I dont like my bicycle)
```

### delete an example sentence

**GraphQL:**
```graphql
mutation deleteSentence {
  deleteSentence(
    polish: "rower"
    english: "bicycle"
    sentence: "I dont like my bicycle."
  )
}
```

**Client:**
```
DELETE_SENTENCE rower bicycle (I dont like my bicycle)
```

### Delete english translation

**Note:** If this is the last translation for a Polish word, the Polish word will also be deleted.

**GraphQL:**
```graphql
mutation deleteTranslation {
  deleteTranslation(
    polish: "rower"
    english: "bicycle"
  )
}
```

**Client:**
```
DELETE_TRANSLATION rower bicycle
```

### Delete polish word with it's translations and example sentences

**GraphQL:**
```graphql
mutation deleteWord {
  deleteWord(
    polish: "rower"
  )
}
```

**Client:**
```
DELETE rower
```

### Update example sentence

**GraphQL:**
```graphql
mutation updateSentence {
  updateSentence(
    polish: "rower"
    english: "bicycle"
    sentence: "I dont like my bicycle."
    newSentence: "I love my bicycle"
  )
}
```

**Client:**
```
UPDATE_SENTENCE rower bicycle (I dont like my bicycle) (I love my bicycle)
```

### Update english translation

**GraphQL:**
```graphql
mutation updateTranslation {
  updateTranslation(
    polish: "rower"
    english: "biccle"
    newEnglish: "bicycle"
  )
}
```

**Client:**
```
UPDATE_TRANSLATION rower biccle bicycle
```

### Update polish word

**GraphQL:**
```graphql
mutation updateWord {
  updateWord(
    polish: "rwer"
    newPolish: "rower"
  )
}
```

**Client:**
```
UPDATE rwer rower
```

### Query word with it's translations and rxamples

**GraphQL:**
```graphql
query select {
  selectWord(polish: "rower") {
    polish
    translations {
      english
      sentences {
        sentence
      }
    }
  }
}
```

**Client:**
```
SELECT rower
```
