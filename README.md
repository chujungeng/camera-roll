
# Camera-Roll CMS

Camera-Roll is a headless CMS for photographer portfolios.
It focuses on organizing the artwork in an intuitive way. 
Pictures are categorized into arbitrary amount of user defined tags, 
and the CMS provides APIs for sourcing data under each tag. 
It also supports customizable albums, 
which photographers may find useful to group together pictures from the same project.

Again, this CMS is built for photographers.
There is no WordPress jargon such as posts and comments.
They don't make sense in the world of aperture, shutter speed and ISO.

## Features

- headless CMS with RESTful APIs written in Go
- image hosting that supports thumbnails, tags and albums
- JSON Web Token(JWT) based access control
- OAuth for admin access, no password required

## Dependencies

- MySQL 8.0
- Google OAuth API

## Usage

Fill in the `config.json` file with your environment info.
Presumably you already have a MySQL database and Google OAuth API set up. 

Compile the code:  
`make`  
This will create a new directory `./bin` with the application binary in it.  

Run the application:  
`cd ./bin`  
`./cameraroll`  

## Example

Here's a sample ReactJS frontend: [https://photography.chujungeng.com/](https://photography.chujungeng.com/)  

In case the above website was censored by your beloved motherland, here're some screenshots of it (you won't be able to see the screenshots either if my entire domain was censored):  

Photos with tag "travel":  
![screenshot of sample frontend](https://chujungeng.com/cameraroll/assets/7b6581cc-2794-48c7-a879-12ea20f246df.jpg)  

Photos with tag "portrait":  
![screenshot of sample frontend](https://chujungeng.com/cameraroll/assets/4e3da456-7bc8-403e-8cc2-30da9795162a.jpg)  

## API Endpoints

GET /images  
get all images  

POST /admin/images  
upload an image  

GET /images/{imageID}  
get the image with id  

GET /images/{imageID}/albums  
get all the albums this image belongs to  

GET /images/{imageID}/tags  
get all the tags this image belongs to  

PUT /admin/images/{imageID}  
modify image with id  

DELETE /admin/images/{imageID}  
delete image with id  

GET /tags  
list all tags  

POST /admin/tags  
add a new tag  

PUT /admin/tags/{tagID}  
modify tag with id  

DELETE /admin/tags/{tagID}  
delete tag with id  

GET /albums  
retrieve all albums  

POST /admin/albums  
add a new album with no pictures in it  

GET /albums/{albumID}  
get the album with albumID  

PUT /admin/albums/{albumID}  
modify album info  

DELETE /admin/albums/{albumID}  
remove album  

GET /albums/{albumID}/images  
get all images from an album  

GET /albums/{albumID}/tags  
get all the tags this album belongs to  

POST /admin/albumImages  
add an image to the album  

DELETE /admin/albums/{albumID}/images/{imageID}  
remove a picture from the album  

GET /tags/{tagID}/albums  
get all albums under the tag  

POST /admin/albumTags  
add tag to the albums  

DELETE /admin/tags/{tagID}/albums/{albumID}  
remove the tag from an album  

GET /tags/{tagID}/images  
get all images under the tag with tagID  

POST /admin/imageTags  
add tag to image  

DELETE /admin/tags/{tagID}/images/{imageID}  
remove the tag from the image  

POST /token/google  
verifies an GoogleID token and responds with an admin JWT if the GoogleID matches admin's  
GoogleID token could be obtained from frontend's OAuth flow  
