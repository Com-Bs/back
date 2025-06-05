#!/bin/bash

echo "Initializing MongoDB with C programming problems..."

# Load environment variables
if [ -f "../.env" ]; then
    export $(cat ../.env | grep -v '#' | xargs)
elif [ -f ".env" ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Connect to MongoDB with authentication and insert problems
mongosh "mongodb://${MONGO_ROOT_USERNAME}:${MONGO_ROOT_PASSWORD}@mongodb:27017/compis?authSource=admin" << 'EOF'

// Drop existing collection
db.problems.drop()


// Insert Add Two Numbers problem
db.problems.insertOne({
  title: "Add Two Numbers",
  description: "Write a C minus function that takes two integers and returns their sum.\n\n## Example\n```\nInput: 5 3\nOutput: 8\n```",
  difficulty: "easy",
  test_cases: [
    { input: "5 3", output: "8" },
    { input: "10 20", output: "30" },
    { input: "-5 15", output: "10" }
  ],
  function_name: "addTwo",
  arguments: [
    { name: "a", type: "int" },
    { name: "b", type: "int" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Even or Odd problem
db.problems.insertOne({
  title: "Even or Odd",
  description: "Write a C minus function that checks if a given integer is even or odd. Return 1 for odd numbers and 0 for even numbers.\n\n## Example\n```\nInput: 4\nOutput: 0\n\nInput: 7\nOutput: 1\n```",
  difficulty: "easy",
  test_cases: [
    { input: "4", output: "0" },
    { input: "7", output: "1" },
    { input: "0", output: "0" },
    { input: "1", output: "1" }
  ],
  function_name: "isOdd",
  arguments: [
    { name: "num", type: "int" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Maximum of Three Numbers problem
db.problems.insertOne({
  title: "Maximum of Three Numbers",
  description: "Write a C minus function that takes three integers and returns the maximum among them.\n\n## Example\n```\nInput: 5 12 8\nOutput: 12\n```",
  difficulty: "easy",
  test_cases: [
    { input: "5 12 8", output: "12" },
    { input: "15 7 9", output: "15" },
    { input: "3 8 8", output: "8" },
    { input: "-1 -5 -3", output: "-1" }
  ],
  function_name: "maxOfThree",
  arguments: [
    { name: "a", type: "int" },
    { name: "b", type: "int" },
    { name: "c", type: "int" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Factorial problem
db.problems.insertOne({
  title: "Factorial",
  description: "Write a C minus function that calculates the factorial of a given positive integer.\n\n## Example\n```\nInput: 5\nOutput: 120\n\nInput: 0\nOutput: 1\n```\n\nFactorial of n (n!) = n × (n-1) × (n-2) × ... × 1",
  difficulty: "medium",
  test_cases: [
    { input: "5", output: "120" },
    { input: "0", output: "1" },
    { input: "1", output: "1" },
    { input: "4", output: "24" }
  ],
  function_name: "factorial",
  arguments: [
    { name: "n", type: "int" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Sum of Natural Numbers problem
db.problems.insertOne({
  title: "Sum of Natural Numbers",
  description: "Write a C minus function that calculates the sum of first n natural numbers.\n\n## Example\n```\nInput: 5\nOutput: 15\n```\n\nSum = 1 + 2 + 3 + 4 + 5 = 15",
  difficulty: "easy",
  test_cases: [
    { input: "5", output: "15" },
    { input: "10", output: "55" },
    { input: "1", output: "1" },
    { input: "100", output: "5050" }
  ],
  function_name: "sumNaturals",
  arguments: [
    { name: "n", type: "int" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Prime Number Check problem
db.problems.insertOne({
  title: "Prime Number Check",
  description: "Write a C minus function that checks if a given number is prime or not. Return 1 if the number is prime, 0 otherwise.\n\n## Example\n```\nInput: 7\nOutput: 1\n\nInput: 8\nOutput: 0\n```\n\nA prime number is a number greater than 1 that has no positive divisors other than 1 and itself.",
  difficulty: "medium",
  test_cases: [
    { input: "7", output: "1" },
    { input: "8", output: "0" },
    { input: "2", output: "1" },
    { input: "1", output: "0" },
    { input: "17", output: "1" }
  ],
  function_name: "isPrime",
  arguments: [
    { name: "num", type: "int" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

// Insert Count Digits problem
db.problems.insertOne({
  title: "Count Digits",
  description: "Write a C minus function that counts the number of digits in a given integer.\n\n## Example\n```\nInput: 12345\nOutput: 5\n\nInput: 7\nOutput: 1\n```",
  difficulty: "easy",
  test_cases: [
    { input: "12345", output: "5" },
    { input: "7", output: "1" },
    { input: "0", output: "1" },
    { input: "999", output: "3" }
  ],
  function_name: "countDigits",
  arguments: [
    { name: "num", type: "int" }
  ],
  created_at: new Date(),
  updated_at: new Date()
})

print("Successfully inserted 7 C minus programming problems!")
print("Collection count:", db.problems.countDocuments())

EOF

echo "MongoDB initialization complete!"