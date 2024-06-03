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

## Project elements
