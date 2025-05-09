FROM zopdev/static-server:latest

# Copy static files from the builder stage
COPY  ./dist website

# COPY ./src/app/favicon.ico /static/favicon.ico

# Expose the port server is running on
EXPOSE 8000

# Start go server
CMD ["/main"]
