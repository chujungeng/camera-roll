import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus} from '@fortawesome/free-solid-svg-icons';

export default function Filmstrip(props) {
    return (
        <Col className="justify-content-center">
        <Row className="align-items-center p-2">
            {
                props.danger? 
                null:
                <Card key="addNew" border="dark" bg="light" text="light" className="mx-2 my-1" style={{ width: '16rem', height: '14rem' }} >
                    <Card.Body as={Button} variant="light" onClick={props.handleAddNew} className="p-4" style={{ textDecoration: 'none', color: 'black' }}>            
                        <Card.Text>
                        <FontAwesomeIcon icon={faPlus} size="3x" />
                        </Card.Text>
                    </Card.Body>
                </Card>
            }
            
            {
                props.photos.map((img) => (
                    <Card key={img.id} border={props.danger? "danger": "dark"} bg="dark" text="light" className="m-2" style={{ width: '16rem', height: '15rem' }}>
                        <Card.Body onClick={img.handleOnClick} as={Button} variant="dark" size="sm" className="p-0">
                        <Card.Img variant="top" src={img.src} style={{width: '100%', height: '14rem', objectFit: 'cover'}} />
                        </Card.Body>
                    </Card>
                ))
            }
        </Row>
        </Col>
    );
}