
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
- npm

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

GET /api/images  
get all images  

POST /api/admin/images  
upload an image  

GET /api/images/{imageID}  
get the image with id  

GET /api/images/{imageID}/albums  
get all the albums this image belongs to  

GET /api/images/{imageID}/tags  
get all the tags this image belongs to  

PUT /api/admin/images/{imageID}  
modify image with id  

DELETE /api/admin/images/{imageID}  
delete image with id  

GET /api/tags  
list all tags  

POST /api/admin/tags  
add a new tag  

PUT /api/admin/tags/{tagID}  
modify tag with id  

DELETE /api/admin/tags/{tagID}  
delete tag with id  

GET /api/albums  
retrieve all albums  

POST /api/admin/albums  
add a new album with no pictures in it  

GET /api/albums/{albumID}  
get the album with albumID  

PUT /api/admin/albums/{albumID}  
modify album info  

DELETE /api/admin/albums/{albumID}  
remove album  

GET /api/albums/{albumID}/images  
get all images from an album  

GET /api/albums/{albumID}/tags  
get all the tags this album belongs to  

POST /api/admin/albumImages  
add an image to the album  

DELETE /api/admin/albums/{albumID}/images/{imageID}  
remove a picture from the album  

GET /api/tags/{tagID}/albums  
get all albums under the tag  

POST /api/admin/albumTags  
add tag to the albums  

DELETE /api/admin/tags/{tagID}/albums/{albumID}  
remove the tag from an album  

GET /api/tags/{tagID}/images  
get all images under the tag with tagID  

POST /api/admin/imageTags  
add tag to image  

DELETE /api/admin/tags/{tagID}/images/{imageID}  
remove the tag from the image  

POST /api/token/google  
verifies an GoogleID token and responds with an admin JWT if the GoogleID matches admin's.  
GoogleID token could be obtained from frontend's OAuth flow.  
with the default admin area, there is no need for this endpoint.  
