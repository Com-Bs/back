#!/bin/bash

echo "Initializing MongoDB with C programming problems..."

# Load environment variables
if [ -f "../.env" ]; then
    export $(cat ../.env | grep -v '#' | xargs)
fi

# Connect to MongoDB with authentication and insert problems
mongosh "mongodb://${MONGO_ROOT_USERNAME}:${MONGO_ROOT_PASSWORD}@localhost:${MONGO_PORT}/compis?authSource=admin" << 'EOF'

// Drop existing collection
db.problems.drop()

// Insert Hello World problem
db.problems.insertOne({
  title: "Hello World",
  description: "Write a C program that prints \"Hello, World!\" to the console.\n\n## Example\n```\nOutput: Hello, World!\n```",
  difficulty: "easy",
  hints: [
    "Use the printf() function to output text",
    "Don't forget to include <stdio.h>",
    "Remember to add a newline character \\n"
  ],
  test_cases: [
    { input: "", output: "Hello, World!" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Add Two Numbers problem
db.problems.insertOne({
  title: "Add Two Numbers",
  description: "Write a C program that reads two integers from input and prints their sum.\n\n## Example\n```\nInput: 5 3\nOutput: 8\n```",
  difficulty: "easy",
  hints: [
    "Use scanf() to read input",
    "Use printf() to output the result",
    "Declare integer variables to store the numbers"
  ],
  test_cases: [
    { input: "5 3", output: "8" },
    { input: "10 20", output: "30" },
    { input: "-5 15", output: "10" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Even or Odd problem
db.problems.insertOne({
  title: "Even or Odd",
  description: "Write a C program that reads an integer and determines if it's even or odd.\n\n## Example\n```\nInput: 4\nOutput: Even\n\nInput: 7\nOutput: Odd\n```",
  difficulty: "easy",
  hints: [
    "Use the modulo operator % to check divisibility",
    "An even number gives remainder 0 when divided by 2",
    "Use if-else statement for conditional logic"
  ],
  test_cases: [
    { input: "4", output: "Even" },
    { input: "7", output: "Odd" },
    { input: "0", output: "Even" },
    { input: "1", output: "Odd" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Maximum of Three Numbers problem
db.problems.insertOne({
  title: "Maximum of Three Numbers",
  description: "Write a C program that reads three integers and finds the maximum among them.\n\n## Example\n```\nInput: 5 12 8\nOutput: 12\n```",
  difficulty: "easy",
  hints: [
    "Compare the first two numbers, then compare the result with the third",
    "Use nested if-else statements",
    "You can also use the ternary operator for concise code"
  ],
  test_cases: [
    { input: "5 12 8", output: "12" },
    { input: "15 7 9", output: "15" },
    { input: "3 8 8", output: "8" },
    { input: "-1 -5 -3", output: "-1" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Factorial problem
db.problems.insertOne({
  title: "Factorial",
  description: "Write a C program that calculates the factorial of a given positive integer.\n\n## Example\n```\nInput: 5\nOutput: 120\n\nInput: 0\nOutput: 1\n```\n\nFactorial of n (n!) = n × (n-1) × (n-2) × ... × 1",
  difficulty: "medium",
  hints: [
    "Use a loop to multiply numbers from 1 to n",
    "Remember that 0! = 1",
    "You can use either for loop or while loop"
  ],
  test_cases: [
    { input: "5", output: "120" },
    { input: "0", output: "1" },
    { input: "1", output: "1" },
    { input: "4", output: "24" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Sum of Natural Numbers problem
db.problems.insertOne({
  title: "Sum of Natural Numbers",
  description: "Write a C program that calculates the sum of first n natural numbers.\n\n## Example\n```\nInput: 5\nOutput: 15\n```\n\nSum = 1 + 2 + 3 + 4 + 5 = 15",
  difficulty: "easy",
  hints: [
    "Use a loop to iterate from 1 to n",
    "Keep adding each number to a sum variable",
    "You can also use the formula n*(n+1)/2"
  ],
  test_cases: [
    { input: "5", output: "15" },
    { input: "10", output: "55" },
    { input: "1", output: "1" },
    { input: "100", output: "5050" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Prime Number Check problem
db.problems.insertOne({
  title: "Prime Number Check",
  description: "Write a C program that checks if a given number is prime or not.\n\n## Example\n```\nInput: 7\nOutput: Prime\n\nInput: 8\nOutput: Not Prime\n```\n\nA prime number is a number greater than 1 that has no positive divisors other than 1 and itself.",
  difficulty: "medium",
  hints: [
    "Check divisibility from 2 to sqrt(n)",
    "If any number divides n evenly, it's not prime",
    "Handle special cases: numbers <= 1 are not prime"
  ],
  test_cases: [
    { input: "7", output: "Prime" },
    { input: "8", output: "Not Prime" },
    { input: "2", output: "Prime" },
    { input: "1", output: "Not Prime" },
    { input: "17", output: "Prime" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Count Digits problem
db.problems.insertOne({
  title: "Count Digits",
  description: "Write a C program that counts the number of digits in a given integer.\n\n## Example\n```\nInput: 12345\nOutput: 5\n\nInput: 7\nOutput: 1\n```",
  difficulty: "easy",
  hints: [
    "Keep dividing the number by 10 until it becomes 0",
    "Count how many divisions you performed",
    "Handle the special case of 0 (it has 1 digit)"
  ],
  test_cases: [
    { input: "12345", output: "5" },
    { input: "7", output: "1" },
    { input: "0", output: "1" },
    { input: "999", output: "3" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

print("Successfully inserted 8 C programming problems!")
print("Collection count:", db.problems.countDocuments())

EOF

echo "MongoDB initialization complete!"