import './App.scss';
import React, { useEffect } from 'react';
import { 
    Link, 
    NavLink, 
    useRoutes,
} from "react-router-dom";
import axios from 'axios';
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import { useSelector, useDispatch } from 'react-redux';
import { logOut, logIn } from '../features/auth/authSlice';
import routes from './routes';

function Navigation() {

    return (
        <Navbar bg="dark" variant="dark" expand="sm" fixed="top">
            <Container className="flex-sm-column">
                <Navbar.Brand className="mr-0" as={Link} to="/"><h1 style={{fontFamily: "'Domine', serif"}}>CAMERA ROLL</h1></Navbar.Brand>
                <Navbar.Toggle aria-controls="basic-navbar-nav" />
                <Navbar.Collapse id="basic-navbar-nav">
                <Nav className="me-auto">
                    <Nav.Item>
                        <Nav.Link as={NavLink} to="/images" key="/images" style={{textTransform: "capitalize"}}>
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
    const apiServer = useSelector((state) => state.api.root);
    const dispatch = useDispatch();
    
    useEffect(() => {
        axios.get(`${apiServer}verify`)
            .then(res => {
                dispatch(logIn());
            }).catch(function (error) {
                dispatch(logOut());
            });

    }, [apiServer, dispatch]);

    const isLoggedIn = useSelector((state) => state.auth.loggedIn);
    const appRoutes = useRoutes(routes(isLoggedIn));

    return (
        <div style={{textAlign: 'center', backgroundColor: 'black', minHeight: "81vh"}}>
            <Navigation />
            
            {appRoutes}
        </div>
    );
};
