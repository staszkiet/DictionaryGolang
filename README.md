What do you need to run the app?

Docker with docker-compose

go language (if you want to use the client app)


Application startup:

clone the repository

go into server folder

run docker-compose build
run docker-compose up -d

Queries and mutations examples:

Create polish-english translation:

mutation createWOrd{
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

in client app:

ADD rower bike (My bike is green.) (I like my bike.)

Add a translation to a polish word existing in a dictionary:

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

in client app:

ADD_TRANSLATION rower bicycle (I like my bicycle) (The bicycle is under repair)

Add example sentence to existing translation:

mutation {
  createSentence(
    	polish: "rower"
      english: "bicycle"
      sentence: "I dont like my bicycle."
  )
}

in clinet app:

ADD_SENTENCE rower bicycle (I dont like my bicycle)

Delete an example sentence:

mutation deleteSentence{
  deleteSentence(
    	polish: "rower"
      english: "bicycle"
      sentence: "I dont like my bicycle."
  )
}

in client app:

DELETE_SENTENCE rower bicycle (I dont like my bicycle)

Delete english part of tranlation (NOTE: if the translation the user deletes is the last one attached to given polish word, the polish word gets also deleted)

mutation deleteTranslation{
  deleteTranslation(
    	polish: "rower"
      english: "bicycle"
  )
}

in client app:

DELETE_TRANSLATION rower bicycle

Delete polish word with all its translations:

mutation deleteWord{
  deleteWord(
    	polish: "rower"
  )
}

DELETE rower

Update example sentence:

mutation updateSentence{
  updateSentence(
    	polish: "rower"
      english: "bicycle"
      sentence: "I dont like my bicycle."
    newSentence: "I love my bicycle"
  )
}

in client app:

UPDATE_SENTENCE rower bicycle (I dont like my bicycle) (I love my bicycle)

Update english part of a translation:

mutation updateTranslation{
  updateTranslation(
    	polish: "rower"
      english: "biccle"
	newEnglish: "bicycle"
  )
}

in client app:

UPDATE_TRANSLATION rower biccle bicycle

Update polish part of the translation

mutation updateWord{
  updateWord(
    	polish: "rwer"
      newPolish: "rower"
  )
}

in client app:

UPDATE rwer rower

Get a word and it's translation and sentences:

query select{
  selectWord(polish: "rower"){
    polish
    translations{
      english
      sentences{
        sentence
      }
    }
  }
}

in client app:

SELECT rower