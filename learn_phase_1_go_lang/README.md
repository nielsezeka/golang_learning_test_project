# Go Parsing Project 

**Author:** na ngo

This project is a Go-based web crawler. It crawls all data from the site https://www.themealdb.com/ and saves the results as JSON files in the output directory.

**The target of this project is:**
- Learn simple Go
- Learn basic libraries
- Learn how to structure a Go project

## Project Structure

- `main.go`: Entry point. Calls the main and detail parsing functions.
- `parsing/main_parse/main_parse.go`: Contains the main page parsing logic.
- `parsing/detail_parse/detail_parse.go`: Contains the detail page parsing logic.
- `output/`: Stores output files (`items.json`, `final_items.json`).

## How to Run

1. Make sure you have Go installed. You can check by running:
   ```sh
   go version
   ```
2. Run the project:
   ```sh
   go run main.go
   ```

The parsed results will be saved in the `output/` directory as JSON files.

---

## How to Run

To run this project, follow these steps:

1. Install Go from https://golang.org/dl/ if you haven't already.
2. Download or clone this project to your computer.
3. Open a terminal and navigate to the project folder.
4. Run the following command:
   ```sh
   go run main.go
   ```
5. Check the `output/` folder for the generated JSON files.

---

## Contact

If you have any questions, please contact na.ngo.
