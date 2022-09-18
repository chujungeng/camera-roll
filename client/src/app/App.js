import './App.scss';
import React from 'react';
import { 
    Link, 
    NavLink, 
    Routes,
    Route } from "react-router-dom";
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import Images from '../routes/Images';
import ImagesIndex from '../routes/ImagesIndex';
import Image from '../routes/Image';
import Albums from '../routes/Albums';
import AlbumsIndex from '../routes/AlbumsIndex';
import Album from '../routes/Album';
import Tags from '../routes/Tags';
import TagsIndex from '../routes/TagsIndex';
import Tag from '../routes/Tag';
import Welcome from '../common/Welcome';
import Login from '../routes/Login';

function Navigation() {

    return (
        <Navbar bg="dark" variant="dark" expand="sm" fixed="top">
            <Container className="flex-sm-column">
                <Navbar.Brand className="mr-0" as={Link} to="/"><h1 style={{fontFamily: "'Domine', serif"}}>CAMERA ROLL</h1></Navbar.Brand>
                <Navbar.Toggle aria-controls="basic-navbar-nav" />
                <Navbar.Collapse id="basic-navbar-nav">
                <Nav className="me-auto">
                    <Nav.Item>
                        <Nav.Link as={NavLink} to="/" key="/images" style={{textTransform: "capitalize"}}>
                            images
                        </Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link as={NavLink} to="/albums" key="/albums" style={{textTransform: "capitalize"}}>
                            albums
                        </Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link as={NavLink} to="/tags" key="/tags" style={{textTransform: "capitalize"}}>
                            tags
                        </Nav.Link>
                    </Nav.Item>
                </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    );
}

export default function App() {
    return (
        <div style={{textAlign: 'center', backgroundColor: 'black', minHeight: "81vh"}}>
            <Navigation />
            <Routes>
                <Route path="/" element={<Images />} >
                    <Route index element={<ImagesIndex />}/>
                    <Route path=":imageID" element={<Image />} />
                </Route>
                <Route path="albums" element={<Albums />}>
                    <Route index element={<AlbumsIndex />}/>
                    <Route path=":albumID" element={<Album />} />
                </Route>
                <Route path="tags" element={<Tags />} >
                    <Route index element={<TagsIndex />}/>
                    <Route path=":tagID" element={<Tag />} />
                </Route>
                <Route
                    path="*"
                    element={
                        <main style={{ padding: "1rem", minHeight: "50vh" }}>
                        <h2 className="text-light">404</h2>
                        <p className="text-light">There's nothing here!</p>
                        </main>
                    }
                />
            </Routes>
            
        </div>
    );
};
