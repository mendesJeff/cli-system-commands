# Use an official GoLang runtime as the base image
FROM golang:1.18.1

# Set the working directory inside the container
WORKDIR /app

# Copy the source code from the current directory to the container
COPY . .

# Build the GoLang application
RUN go build -o main .

# Expose the port on which your GoLang application listens
#EXPOSE 5000

# Define the command to run your GoLang application
CMD ["./main"]

