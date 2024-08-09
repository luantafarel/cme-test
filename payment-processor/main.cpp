#include <iostream>
#include <string>

const int MAX_ACCOUNTS = 10;

struct Account {
    std::string accountId;
    std::string ownerName;
    double balance;
};

struct Transaction {
    std::string transactionId;
    std::string fromAccountId;
    std::string toAccountId;
    double amount;
};

Account accounts[MAX_ACCOUNTS];
Transaction transactions[MAX_ACCOUNTS];
int accountCount = 0;
int transactionCount = 0;

std::string generateUniqueId() {
    return std::to_string(rand() % 1000000);
}

int findAccountIndex(const std::string& accountId) {
    for (int i = 0; i < accountCount; i++) {
        if (accounts[i].accountId == accountId) {
            return i;
        }
    }
    return -1;
}

std::string createAccount(const std::string& ownerName, double initialBalance) {
    if (accountCount >= MAX_ACCOUNTS) {
        std::cerr << "Cannot create more accounts. Limit reached." << std::endl;
        return "";
    }

    accounts[accountCount].accountId = generateUniqueId();
    accounts[accountCount].ownerName = ownerName;
    accounts[accountCount].balance = initialBalance;

    std::cout << "Account created for " << ownerName << " with ID: " << accounts[accountCount].accountId 
              << " and initial balance: $" << initialBalance << std::endl;

    accountCount++;
    return accounts[accountCount - 1].accountId;
}

bool processTransaction(const std::string& fromAccountId, const std::string& toAccountId, double amount) {
    int fromIndex = findAccountIndex(fromAccountId);
    int toIndex = findAccountIndex(toAccountId);

    if (fromIndex == -1 || toIndex == -1) {
        std::cerr << "Invalid account ID(s)." << std::endl;
        return false;
    }

    if (accounts[fromIndex].balance < amount) {
        std::cerr << "Insufficient funds in the source account." << std::endl;
        return false;
    }

    // Print balances before the transaction
    std::cout << "\nBefore Transaction:" << std::endl;
    std::cout << accounts[fromIndex].ownerName << "'s balance: $" << accounts[fromIndex].balance << std::endl;
    std::cout << accounts[toIndex].ownerName << "'s balance: $" << accounts[toIndex].balance << std::endl;

    // Process the transaction
    accounts[fromIndex].balance -= amount;
    accounts[toIndex].balance += amount;

    transactions[transactionCount].transactionId = generateUniqueId();
    transactions[transactionCount].fromAccountId = fromAccountId;
    transactions[transactionCount].toAccountId = toAccountId;
    transactions[transactionCount].amount = amount;
    transactionCount++;

    // Print transaction details
    std::cout << "\nTransaction Details:" << std::endl;
    std::cout << "Transaction ID: " << transactions[transactionCount - 1].transactionId << std::endl;
    std::cout << "Amount $" << amount << " transferred from " << accounts[fromIndex].ownerName 
              << " (ID: " << fromAccountId << ") to " << accounts[toIndex].ownerName 
              << " (ID: " << toAccountId << ")." << std::endl;

    // Print balances after the transaction
    std::cout << "\nAfter Transaction:" << std::endl;
    std::cout << accounts[fromIndex].ownerName << "'s balance: $" << accounts[fromIndex].balance << std::endl;
    std::cout << accounts[toIndex].ownerName << "'s balance: $" << accounts[toIndex].balance << std::endl;

    return true;
}

int main() {
    // Create two new accounts
    std::string aliceId = createAccount("Alice", 500.0);
    std::string bobId = createAccount("Bob", 300.0);

    // Process a transaction using the correct account IDs
    bool success = processTransaction(aliceId, bobId, 150.0);

    // Print the result of the transaction
    if (success) {
        std::cout << "\nTransaction was successful." << std::endl;
    } else {
        std::cout << "\nTransaction failed." << std::endl;
    }

    return 0;
}
