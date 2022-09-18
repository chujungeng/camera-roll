import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus} from '@fortawesome/free-solid-svg-icons';
import { useSelector, useDispatch } from 'react-redux';
import { logOut } from '../features/auth/authSlice';

export default function TagsIndex(props) {
    const [tags, setTags] = useState([]);
    const apiServer = useSelector((state) => state.api.root);

    const axiosInstance = axios.create();
    const dispatch = useDispatch();

    axiosInstance.interceptors.response.use(
        (response) => response,
            (error) => {
                if (error.response.status === 401) {
                    dispatch(logOut());
                    return error.response;
                }
    
                return Promise.reject(error);
            }
    )

    useEffect(() => {
        axiosInstance.get(`${apiServer}tags`)
            .then(res => {
                if (Array.isArray(res.data) && res.data.length !== 0) {
                    setTags(res.data);
                }
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });
    }, [apiServer, axiosInstance]);

    return (
        <main style={{flex: 1, minHeight: "55vh"}}>
        <Container>
        <Col className="justify-content-center">
            <Row className="align-items-center p-2">
            {
                tags.map((tag) => (
                    <Card key={tag.id} border="light" bg="dark" text="light" className="mx-2 my-1" style={{ width: '16rem', height: '10rem' }}>
                        <Card.Body as={Link} style={{ textDecoration: 'none' }} to={`/tags/${tag.id}`}>
                            <Card.Title className="text-light" style={{textTransform: "capitalize"}}>
                                {tag.name}
                            </Card.Title>
                        </Card.Body>
                        <Card.Body>
                            <Card.Link className="text-secondary" href="#">Edit</Card.Link>
                            <Card.Link className="text-danger" href="#">Delete</Card.Link>
                        </Card.Body>
                    </Card>
                ))
            }
                <Card key="addNew" border="light" className="mx-2 my-1" style={{ width: '16rem' }} >
                    <Card.Body as={Button} variant="light" onClick={props.handleAddNew} className="p-4" style={{ width: '16rem', height: '10rem' }}>            
                        <Card.Text>
                        <FontAwesomeIcon icon={faPlus} size="3x" />
                        </Card.Text>
                    </Card.Body>
                </Card>
            </Row>
        </Col>
        </Container>
        </main>
    );
}