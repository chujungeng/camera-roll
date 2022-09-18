import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from "react-router-dom";
import axios from 'axios';
import InfiniteScroll from 'react-infinite-scroll-component';
import Filmstrip from '../common/Filmstrip';
import Toggle from 'react-toggle';
import 'react-toggle/style.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { useSelector, useDispatch } from 'react-redux';
import { logOut } from '../features/auth/authSlice';

export default function Images() {
    const [images, setImages] = useState([]);
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);
    const [deletion, setDeletion] = useState(false);
    const apiServer = useSelector((state) => state.api.root);
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

        axiosInstance.get(`${apiServer}images`)
            .then(res => {
                if (Array.isArray(res.data) && res.data.length !== 0) {
                    setImages(res.data);
                }
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });
    }, [apiServer, dispatch]);

    const fetchMoreData = useCallback(() => {
        const nextPage = page + 1;
        axios.get(`${apiServer}images?page=${nextPage}`)
            .then(res => {
                if (Array.isArray(res.data) && res.data.length !== 0) {
                    setImages(prevImages => prevImages.concat(res.data));
                    setPage(nextPage);
                } else {
                    setHasMore(false);
                }
            }).catch(function (error) {
                setHasMore(false);
            });
    }, [page, apiServer]);

    const toggleDelete = useCallback(() => {
        setDeletion(curr => !curr);
    }, []);

    // TODO: deletion route
    const photos = images.map(img => ({
        id: img.id,
        src: img.thumbnail,
        handleOnClick: () => {navigate(`/images/${img.id}`, {replace: false})},
        handleDelete: () => {
            axios.delete(`${apiServer}images/${img.id}`).then(setImages(prevImages => prevImages.filter(i => i.id != img.id)));
        },
    }));

    return (
        <main style={{flex: 1, minHeight: "55vh"}}>
            <Container>
                <Row className="align-items-center p-1">
                    <Col className="d-flex justify-content-start">
                        <span className="text-light mr-2">Deletion Mode</span>
                        <Toggle
                            defaultChecked={false}
                            icons={false}
                            onChange={toggleDelete} />
                    </Col>
                </Row>
            

                <InfiniteScroll
                    dataLength={images.length}
                    next={fetchMoreData}
                    hasMore={hasMore}
                >
                    <Filmstrip 
                        photos={photos}
                        enableSelection={false}
                        onClick={(idx) => {deletion? photos[idx].handleDelete(): photos[idx].handleOnClick()}}
                    />
                </InfiniteScroll>
            </Container>
        </main>
    );
}