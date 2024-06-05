# Expenses app

This is the repository for my Expenses app

## Project goal

### Bulding an application capable of managing my personal expenses and budget

For a couple of years now I have been using a Google Sheets spreadsheet to manage my personal finances. This has went as expected. In the beggining the spreadsheet was simple and did the job, but as time went by it became more and more complex where as now I have extremely long formulas that are completely unmaintanable.

This application will have to be able to replace my spreadsheet. To acomplish this it will have to provide a couple of base functionalities but also to allow for chaning design specifications. I have never stoped to tinker with my budget spreadsheet and if the application is too rigidly built and doesn't allow for extentions, it won't be a good replacement.

### Learning the basics of Go

My main language is javascript and I am experienced with JS, TS, Node.js, React, Next JS and others. That's all well and good but a programmer that is stuck on a single language or framework runs the risk of falling behind. I feel that by learning Go I will extend my abilities as a programmer and will widen my horizon to allow me to become a better JS developer.

### Trying some unorthodox designs

This project will feature some weird designs. SQLite with one DB per user? Custom interpreted language to analize the data? Making the DB module horizontally scalable?

Some of the decisions I take in this project will not be particularly good but I will take them nontheless because they are fun and will provide for an opportunity to learn.

## Key concepts and specifications

### Accounts

An **Account** is a place where the user stores money. It can be a bank account, the cash he has on hand or even as simle as a piggy bank. Having your accounts represented in the application helps with tracking what money is where.

### Categories

By **Category** the app means a budget category. The budget **category** helps the user plan how much does he want to spend for different things. It also helps him track what does he spend his money on.

### Expenses

**Expenses** are individual sums of money that the user has spent on something. The expense tracks a couple of key data:

- Amount - how much has the user spent
- From which account has the user taken the money from
- From which category does the exoense come out of
- Tags - additional description tags that allow for more fine grained tracking of individual expenses

### Basic workflow

1. Creating an user profile
   When the user creates a profile (or has a profile created for him) he will be expected to create Accounts for each of his payment accounts.
   He will also be expected to create budget categories, by listing:

   - Category name
   - The Period over which he expects to input funds. e.g.: One a month when he receives his salary or once a year when he receives some sort of yearly payment
   - The amount he expect to add in the beggining of each perid
   - The amount he wants to spend over the period, or the spending limit

2. Adding money for the first time
   When money is added, it goes in a specific account and is countet as **_Free funds_**.
   After the user has added the free funds that are in each of his accounts, he can reset the budget categories for the first time. On this initial reset he will input how much money is supposed to be in each of the categories.

3. Adding expenses
   After the accounts and categories are setup the user can start tracking expenses. Each expense has to have at least one tag so after a time the user can know why he made it.

4. Inputing more money
   When a category period end, the user can add more money to the free funds and then transfer that money from the free funds to one or more budget categories. The budget cateogry input amount and spending limit can be changed at this point.
   When he adds money, the category period gets reset and an archive record is created for the period, for future reference and statistics purposes.

## Project architecture

### Stage 1

The project is a golang server that keeps the sqlite databases localy and doesn't communicate with external services.

### Stage 2

For Stage 2 I will split the project in to two parts. The main server, that renders content for the users and a database server, that contains all of the logic for connecting and communicating with the db. This seems like an unnecessary transition but it's a neccessary step if I want to make a horizontally scalable sqlite database. It also decouples the user facing part from the buisness logic and I may end up converting that to NextJS. The communication between services will happen in gRPC firstly because it's a fun new thing to try, seondly because it's a flexible and robust way to organize service-to-service communication.

## Project elements

### Stage 1

At the time of writing this document, a couple of design decisions where made. The project is in to active testing to iterate the design and remove bugs.

#### Database

The proect uses SQLite. It has a DB per user. The DB has a User table that stores the hashed user password and also tables for all other data.
When a user loges in the server opens a db connection specifically for the user. At this moment there is no way to create a new user, appart from running the app with a -seed flag, which creates a test user with some data setup. The db file is stored in the /db folder.

One DB per user is fine option for this type of project, because the users have no interaction between them. It comes with a couple of benefits:

- The queries get simplified. That is not a huge deal, but it does lead to a couple of table design improvemnts.
- It reduces the risk for a user to get access to someone else's data
- It makes it simple to give a user a copy of his data or to delete his data if he requests it
- It opens the option for horizontaly scaling the database. Again not really important for this project, but still fun
- It allows for simpler migrations of the db schema. When there is only one db, migrations have to happen at once for everyone. In this db-per-user architecture migrations can be applyed on an individual basis. This allows for testing the migration on a subset of users to see for bugs or to performe the migration at a convinient for the user time.
- It may have some performance benefits in the case of a large user base. Each db write starts a transaction that usually locks a couple of tables. This could become a bottleneck when using a single db for everyone.

I decided to put as much of the buisiness logic as possible in the database for a couple of reasons. Firstly beacause I don't have a ton of experience with SQL and found it fun. Secondly because this opens a possibility to use the db in a mobile application and to keep it on device for offline access. In this scenario it would be usefull to have most of the buisiness logic in the db and not to have to rewirte it on a while other language.

I decided to keep some redundand data. For example I am storing the amount a user has left to spend in a category, as a column. This of course can be calculated by taking the inital amount in the category and subtracting the expenses from the start of the period. What I don't like about this solution is that you have to query the expenses by unidexed fields, once for every category. By keeping the value as a column I overcome that. The downside is that I have to recalculate this value every time the user adds an expense. This is the point where using SQLite became something of a chalange. In a more capable db this could be done through a database procedure but SQLite doesn't have them.

I implemented something similar to DB procedure through a combination of views and triggers. For every "procedure" I create a View that exposes the required columns. I then write a trigger that runs on the INSTEAD OF event and performes the necessary actions. This has a couple of drawbacks but seems to work fine. **_Important:_** an uninforced rule is that the app will comunicate with the db only through these simulated procedures and will not make direct calls to the main tables, in order to protect the integrity of the data.

The main drawback of these procedures is that the app can't fetch the last inserted id. As a workaround I configure the db so the transaction lock is exclusive. This means that the last inserted row in a column will always be the larges id. Not ideal but it get's the job done.

#### Go Server and routing

The web part of the project has an app config that is passed as a dependency to other modules. It contains basic app configuration and access to loggers
The routing is achieved through a Chi router.

#### Handlers

The handler module takes the app config and a map of active db connections as a dependency injections.
It provides helper methods for managing the user Session contents and adding data to it.
It provides handlers for the different endpoints.

The default behaviour is that each page has a GET handler that renders some content.
When the user submits data through a form, the data is sent to a POST handler.
When the POST handler does it's thing, it redirects the user to the appropriate page and inserts data in the user session. The data can be a success or error message and also it can be form data.
Each form is stored in a map and has a specific key associated with it.

#### Page rendering

The server uses Templ to render the pages. Each Templ template expects some data to be provided by the handler.

#### Database repository

The handler package expects to have access to DB repositories, one for each user. The DB repository has access to the app config and a DB connection, through dependecy injection. It also provides methods for interacting with the database.

#### Frontend code

I've made a point not to use external libraries. Not for any good reason, but because I find it fun to develop js code from scratch.

### Stage 2

#### DB Controller

The database controller will take care of interacting with the database and authenticating the user.
