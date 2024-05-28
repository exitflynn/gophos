# Cracking Sophos Creds

This is a personal exploration into brute-forcing Sophos login credentials. It's meant for educational purposes and to understand the mechanics behind brute-force attacks. 

### What Does This Do?

This project attempts to brute-force the Sophos login credentials for university students. The underlying idea at the core is that everyone's passwords are picked from a repeating set of finite randomized combinations, and hence are not truly random. The script goes through every username for a given password to find matches.

### Some other info

1. It requires atleast one correct credential pair which are taken from the Environment Variables `SOPHOS_USERNAME` and `SOPHOS_PASSWORD`.

2. Potential passwords are read from `passwords.csv` and successful logins are written to `matched.csv`.

3. It tries logging in with different username and password combinations and logs the results, parsing the XML responses to check the login status.

4.  Iterates through each username and password combination. Sidesteps getting rate-limited by logging in and logging out with the legitimate credentials to reset the state after a few wrong attempts.

### Why Did I Make This?

for the lolz :P

## Getting Started

### Prerequisites

1. **Go**: Make sure you have Go installed. If not, you can download it from [here](https://golang.org/dl/).
2. **Environment Variables**: Set `SOPHOS_USERNAME` and `SOPHOS_PASSWORD` with your correct credentials.
3. Should be **connected to a Sophos Network**.

### Setting Up

1. **Clone the Repository**:

   ```sh
   git clone https://github.com/exitflynn/gophos
   cd gophos
   ```

2. **Prepare CSV Files**:
   - Create a `passwords.csv` file in the root directory with potential passwords.
   - Ensure you have a `matched.csv` file where successful logins will be recorded.

3. **Run the Script**:

   ```sh
   go run main.go
   ```
