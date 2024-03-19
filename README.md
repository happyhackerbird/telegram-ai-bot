# Telegram AI Chatbot  

Telegram bot that allows users to create and chat with multiple personalized AI assistants using OpenAI and Perplexity API. The user can create an assistant profile by setting a custom system prompt for the model and choose between different LLM models. 

Each profile responds based on a relevant history that is implemented with Milvus vector database (through a process called Retrieval Augmented Generation). 

This methods works by using vector embeddings of the message text in order to retrieve the most relevant messages (n=7) via a similarity search on the vector DB. The prompt to the LLM is augmented with this semantic context, hence the name. The memory is persistent between different chat sessions and specific to the user's conversation (chat). 
Previously the bot had a short term memory achieved using a sliding window history. 

Questions: 
- Do the individual messages need to be chunked evenly before storage?
- Refining the context retrieval

## To-do
- [ ] add prompt generation from messages to help users tailor the assistant to their needs
- [ ] (in progress) store the prompts in the db & retrieve them + ability to edit them 
      
## To run your own bot 
1. Create a new bot and get the API token from the Telegram Botfather account
2. Get Perplexity and OpenAI API tokens
3. Create Milvus Cluster (via Zilliz) and get the API token
4. Create the DB schema by running the ```database.Migrate()``` function the first time
5. Run the server supplying the API keys in the .env file
6. Access the chatbot under its Telegram name/link

## Architecture 
Architecture Example I followed: https://github.com/eikoshelev/go-telegram-bot-example  
This is supposed to be clean architecture, I don't fully understand this design and am unsure if the way I implemented is correct under this approach.
