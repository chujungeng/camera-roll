
# Camera-Roll CMS

Camera-Roll is a headless CMS for photographer portfolios.
It focuses on organizing the artwork in an intuitive way. 
Pictures are categorized into arbitrary amount of user defined tags, 
and the CMS provides APIs for sourcing data under each tag. 
It also supports customizable albums, 
which photographers may find useful to group together pictures from the same project.

Again, this CMS is built for photographers.
There is no WordPress jargon such as posts and users which were meant for blogs.
They don't make sense in the world of aperture, shutter speed and ISO.

## Features

- headless CMS with RESTful APIs written in Go
- organizes pictures into albums and tags
- JSON Web Token(JWT) based access control
- OAuth for admin access, no password required (Not yet implemented)

## Dependencies

- MySQL 8.0

## Usage

(Not yet implemented)

## Example

(Not yet implemented)

## API endpoints

/images  
GET: get all images  
POST: upload image  

/images/{imageID}  
GET: get the image with id  
PUT: modify image with id  
DELETE: delete image with id  

/tags  
GET: list all tags  
POST: add a new tag  

/tags/{tagID}  
PUT: modify tag with id  
DELETE: delete tag with id  

/albums  
GET: retrieve all albums  
POST: add a new album with no pictures in it  

/albums/{albumID}  
GET: retrieve an album and its pictures  
PUT: modify album info  
DELETE: remove album  

/albums/{albumID}/images  
GET: get all images from an album  

/albumImages  
POST: add an image to the album  

/albums/{albumID}/images/{imageID}  
DELETE: remove a picture from the album  

/tags/{tagID}/albums  
GET: get all albums under the tag  

/albumTags  
POST: add tag to the albums  

/tags/{tagID}/albums/{albumID}  
DELETE: remove the tag from an album  

/tags/{tagID}/images  
GET: get all images under the tag  

/imageTags  
POST: add tag to image  

/tags/{tagID}/images/{imageID}  
DELETE: remove the tag from an image  
