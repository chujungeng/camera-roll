import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from "react-router-dom";
import InfiniteScroll from 'react-infinite-scroll-component';
import PhotoAlbum from 'react-photo-album';
import Toggle from 'react-toggle';
import 'react-toggle/style.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Lightbox from "yet-another-react-lightbox";
import Fullscreen from "yet-another-react-lightbox/plugins/fullscreen";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import "yet-another-react-lightbox/styles.css";
import "yet-another-react-lightbox/plugins/thumbnails.css";
import { useSelector, useDispatch } from 'react-redux';
import { logOut } from '../features/auth/authSlice';

const views = ["images", "albums"];

export default function Tag(props) {
    const params = useParams();
    const tagID = params.tagID;

    const [tag, setTag] = useState({});
    const [images, setImages] = useState([]);
    const [albums, setAlbums] = useState([]);
    const [activeView, setActiveView] = useState(0);
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);
    const [activeIndex, setActiveIndex] = useState(-1);
    const apiServer = useSelector((state) => state.api.root);
    const navigate = useNavigate();

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
        axiosInstance.get(`${apiServer}tags/${tagID}`)
            .then(res => {
                setTag(res.data);
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });
    }, [apiServer, tagID, navigate, axiosInstance]);

    useEffect(() => {
        if (activeView >= 0 && activeView < views.length) {
            axiosInstance.get(`${apiServer}tags/${tagID}/${views[activeView]}`)
            .then(res => {
                    if (activeView === 0) {
                        setImages(res.data);
                    } else {
                        setAlbums(res.data);
                    }
                    setHasMore(true);
                    setPage(1);
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });
        }
    }, [apiServer, tagID, activeView, axiosInstance]);

    const toggleView = useCallback(() => {
        setActiveView(curr => (curr + 1) % views.length);
    }, []);

    const fetchMoreData = useCallback(() => {
        const nextPage = page + 1;
        let canceled = false;
        axiosInstance.get(`${apiServer}tags/${tagID}/${views[activeView]}?page=${nextPage}`)
            .then(res => {
                if (!canceled) {
                    if (Array.isArray(res.data) && res.data.length !== 0) {
                        if (activeView === 0) {
                            setImages(prev => prev.concat(res.data));
                        } else {
                            setAlbums(prev => prev.concat(res.data));
                        }
                        setPage(nextPage)
                    } else {
                        setHasMore(false);
                    }
                }
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });
        return () => (canceled = true);
    }, [apiServer, tagID, page, activeView, axiosInstance]);

    const photos = activeView === 0? 
        images.map((img) => ({
            imageID: img.id,
            key: img.id,
            src: img.thumbnail,
            width: img.width_thumb,
            height: img.height_thumb,
            title: img.title,
            desc: img.description})): 
        albums.map((alb) => ({
            albumID: alb.id,
            key: alb.id,
            src: alb.cover? alb.cover.thumbnail: '/default.svg',
            width: alb.cover? alb.cover.width_thumb: '300',
            height: alb.cover? alb.cover.height_thumb: '300',
            title: alb.title,
            desc: alb.description,
        }));

    return (
        <main style={{flex: 1, minHeight: "55vh"}}>
            <Container>
                <Row className="align-items-center p-1">
                    <Col className="d-flex justify-content-start">
                        <h2 className="text-light" style={{textTransform: "capitalize"}}>
                            {tag === null? `Tag #${tagID}`: tag.name}
                        </h2>
                    </Col>
                    
                    <Col className="d-flex justify-content-end">
                        <span className="text-light mr-2">AlbumView</span>
                        <Toggle
                            defaultChecked={false}
                            icons={false}
                            onChange={toggleView} />
                    </Col>
                </Row>
            </Container>
            
            <InfiniteScroll
                dataLength={activeView === 0? images.length: albums.length}
                next={fetchMoreData}
                hasMore={hasMore}
            >
                <PhotoAlbum 
                    layout="masonry" 
                    spacing={(containerWidth) => {
                        if (containerWidth >= 1200) return 10;
                        if (containerWidth >= 600) return 6;
                        if (containerWidth >= 300) return 4;
                        return 2;
                    }}
                    columns={(containerWidth) => {
                        if (containerWidth >= 1200) return 4;
                        if (containerWidth >= 600) return 3;
                        if (containerWidth >= 300) return 2;
                        return 1;
                    }}
                    targetRowHeight={(containerWidth) => {
                        if (containerWidth >= 1200) return containerWidth / 4;
                        if (containerWidth >= 600) return containerWidth / 3;
                        if (containerWidth >= 300) return containerWidth / 2;
                        return containerWidth;
                    }}
                    photos={photos}
                    onClick={(event, photo, index) => {
                        setActiveIndex(index);
                    }}
                >
                </PhotoAlbum>
            </InfiniteScroll>

            {
            activeView === 0?    
                <Lightbox
                open={activeIndex >= 0}
                close={() => setActiveIndex(-1)}
                index={activeIndex}
                slides={images.map((img) => ({
                    src: img.path
                }))}
                plugins={[Zoom, Fullscreen, Thumbnails]}
                carousel={{ finite: true, preload: 1 }}
                thumbnails={{ width: 60, height: 40, padding: 2}}
                />:
                <AlbumView
                    apiServer={props.apiServer} 
                    albums={albums} 
                    index={activeIndex} 
                    handleClose={() => setActiveIndex(-1)}
                />
            }
        </main>
    );
}

class AlbumView extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            images: []
        };
    }

    componentDidMount() {
        document.addEventListener('contextmenu', (e) => {
          e.preventDefault();
        });

        if (this.props.albums !== null && this.props.index >= 0 && this.props.index < this.props.albums.length) {
            const album = this.props.albums[this.props.index];
            axios.get(`${this.props.apiServer}albums/${album.id}/images`)
            .then(res => {
                const images = res.data
                this.setState({images});
            }).catch(function (error) {
                if( error.response ){
                    // console.log(error.response.data); // => the response payload 
                }
            });
        }
    };

    componentDidUpdate(prevProps) {
        if (prevProps.index !== this.props.index 
            && this.props.albums !== null 
            && this.props.index >= 0 
            && this.props.index < this.props.albums.length) {
                const album = this.props.albums[this.props.index];

                // clear previous state
                this.setState({images: []});

                // fetch images from backend API
                axios.get(`${this.props.apiServer}albums/${album.id}/images`)
                .then(res => {
                    const images = res.data
                    this.setState({
                        images: images,
                    });
                }).catch(function (error) {
                    if( error.response ){
                        // console.log(error.response.data); // => the response payload 
                    }
                });
        }
    }

    render() {
        if (this.props.albums === null 
            || this.props.index < 0 
            || this.state.images === null 
            || this.state.images.length === 0) {
            return null;
        }

        return (
            <Lightbox
                open={true}
                close={this.props.handleClose}
                slides={this.state.images.map((img) => ({
                    src: img.path
                }))}
                plugins={[Zoom, Fullscreen, Thumbnails]}
                carousel={{ finite: true }}
                thumbnails={{ width: 60, height: 40, padding: 2}}
            />
        );
    } 
}