#!/bin/bash

# Change to the directory with our code that we plan to work from 
cd "$GOPATH/src/github.com/samueldaviddelacruz/lenslocked.com"

echo "==== Releasing lenslocked-project-demo.net ===="
echo " Deleting the local binary if it exists (so it isn't uploaded)..."
rm lenslocked.com
echo " Done!"

echo " Deleting existing code..."
ssh root@lenslocked-project-demo.net "rm -rf /root/go/src/github.com/samueldaviddelacruz/lenslocked.com"
echo " Code deleted successfully!"


echo " Uploading code..."
# The \ at the end of the line tells bash that our
# command isn't done and wraps to the next line.
rsync -avr --exclude '.git/*' --exclude 'tmp/*' \
--exclude 'images/*' ./ \
root@lenslocked-project-demo.net:/root/go/src/github.com/samueldaviddelacruz/lenslocked.com
echo " Code uploaded successfully!"

echo " Go getting deps..."
ssh root@lenslocked-project-demo.net "export GOPATH=/root/go; \
/usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@lenslocked-project-demo.net "export GOPATH=/root/go; \
/usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@lenslocked-project-demo.net "export GOPATH=/root/go; \
/usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@lenslocked-project-demo.net "export GOPATH=/root/go; \
/usr/local/go/bin/go get github.com/lib/pq"
ssh root@lenslocked-project-demo.net "export GOPATH=/root/go; \
/usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@lenslocked-project-demo.net "export GOPATH=/root/go; \
/usr/local/go/bin/go get github.com/gorilla/csrf"

ssh root@lenslocked-project-demo.net "export GOPATH=/root/go; \
/usr/local/go/bin/go get gopkg.in/mailgun/mailgun-go.v3"

echo " Building the code on remote server..."
ssh root@lenslocked-project-demo.net 'export GOPATH=/root/go; \
cd /root/app; \
/usr/local/go/bin/go build -o ./server \
$GOPATH/src/github.com/samueldaviddelacruz/lenslocked.com/*.go'
echo " Code built successfully!"

echo " Moving assets..."
ssh root@lenslocked-project-demo.net "cd /root/app; \
cp -R /root/go/src/github.com/samueldaviddelacruz/lenslocked.com/assets ."
echo " Assets moved successfully!"
echo " Moving views..."
ssh root@lenslocked-project-demo.net "cd /root/app; \
cp -R /root/go/src/github.com/samueldaviddelacruz/lenslocked.com/views ."
echo " Views moved successfully!"
echo " Moving Caddyfile..."
ssh root@lenslocked-project-demo.net "cd /root/app; \
cp /root/go/src/github.com/samueldaviddelacruz/lenslocked.com/Caddyfile ."
echo " Views moved successfully!"


echo " Restarting the server..."
ssh root@lenslocked-project-demo.net "sudo service lenslocked-project-demo restart"
echo " Server restarted successfully!"
echo " Restarting Caddy server..."
ssh root@lenslocked-project-demo.net "sudo service caddy restart"
echo " Caddy restarted successfully!"
echo "==== Done releasing lenslocked-project-demo ===="



