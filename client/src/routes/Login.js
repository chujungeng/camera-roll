import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Button from 'react-bootstrap/Button';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faGoogle } from '@fortawesome/free-brands-svg-icons';

export default function Login() {

    return (       
            <ButtonGroup vertical className="my-4">
                <Button href="/auth/google/login" variant="outline-light">
                    <FontAwesomeIcon icon={faGoogle} className="mr-2"/>Sign In With Google
                </Button>
            </ButtonGroup>
    );
}