# Telegram AI Chatbot  

A Telegram bot that allows users to chat to a personalized AI assistant. The user can set a custom system prompt for the AI model and select the model. Messages are stored in a vector database and retrieved to provide a relevant context for new user input. 

Memory is currently achieved using the Retrieval Augmented Generation (RAG). This methods works by using vector embeddings of the message text in order to retrieve the most relevant messages (n=7) via a similarity search on the vector DB. Database used is Milvus (Cloud). The retrieved message provide the context for the user input, everything is sent along in the API call. 
This memory is persistent between different chat sessions and specific to the user's conversation (chat). 
Previously the bot had a short term memory achieved using the sliding window method. 
Questions: 
- Do the individual messages need to be chunked evenly before storage?
- Refining the context retrieval

## To-do
- [ ] add ability to create multiple threads so that the user can have multiple ongoing chats (make the bot profile persistent for each thread)
- [ ] check if error handling was done correctly everywhere
      
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
