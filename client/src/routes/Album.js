import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from "react-router-dom";
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Filmstrip from '../common/Filmstrip';
import Toggle from 'react-toggle';
import 'react-toggle/style.css';
import { useSelector, useDispatch } from 'react-redux';
import { logOut } from '../features/auth/authSlice';

export default function Album() {
    const params = useParams();
    const albumID = params.albumID;
    const apiServer = useSelector((state) => state.api.root);
    const [album, setAlbum] = useState({})
    const [images, setImages] = useState([]);
    const [tags, setTags] = useState([]);
    const [deletion, setDeletion] = useState(false);
    const navigate = useNavigate();
    const dispatch = useDispatch();


    useEffect(() => {
        const axiosInstance = axios.create();

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

        axiosInstance.get(`${apiServer}albums/${albumID}`)
            .then(res => {
                setAlbum(res.data);
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });

        axiosInstance.get(`${apiServer}albums/${albumID}/images`)
            .then(res => {
                setImages(res.data);
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });;
        
        axiosInstance.get(`${apiServer}albums/${albumID}/tags`)
            .then(res => {
                setTags(res.data);
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });;
    }, [apiServer, albumID, navigate, dispatch]);

    const toggleDelete = useCallback(() => {
        setDeletion(curr => !curr);
    }, []);

    // TODO: deletion route
    const photos = images.map(img => ({
        id: img.id,
        src: img.thumbnail,
        handleOnClick: () => {navigate(`/images/${img.id}`, {replace: false})},
    }));

    return (
        <main style={{flex: 1, minHeight: "55vh"}}>
            <Container>
                <Row className="align-items-center p-1">
                    <Col className="d-flex justify-content-start">
                        <h2 className="text-light" style={{textTransform: "capitalize"}}>
                            {album === null? 
                            `Album #${albumID}`: 
                            album.title? album.title: "Untitled"}
                        </h2>
                    </Col>
                </Row>
                <Row className="align-items-center p-1">
                    <Col className="d-flex justify-content-start">
                        <h4 className="text-light" style={{textTransform: "capitalize"}}>
                            {album === null? 
                            `Album #${albumID}`: 
                            album.description? album.description: "No Description"}
                        </h4>
                    </Col>
                </Row>
                <Row className="align-items-center p-1">
                    <Col className="d-flex justify-content-start">
                        {
                            album === null?
                            null:
                            <span className="text-light">Created At: {album.created_at}</span>
                        }
                    </Col>
                </Row>
                
                {
                    tags.length === 0?
                    null:
                    <Row className="align-items-center p-1">
                        <Col className="d-flex justify-content-start">
                            {
                                tags.map(tag => (<span className="text-light mr-2" key={tag.id}> {tag.name}</span>))
                            }
                        </Col>
                    </Row>
                }

                <Row className="align-items-center p-1">
                    <Col className="d-flex justify-content-start">
                            <span className="text-light mr-2">Deletion Mode</span>
                            <Toggle
                                defaultChecked={false}
                                icons={false}
                                onChange={toggleDelete} />
                    </Col>
                </Row>

                <Filmstrip 
                    photos={photos}
                    danger={deletion}
                    handleAddNew={()=>{console.log('not yet implemented')}}
                />
            </Container>
        </main>
    );
}