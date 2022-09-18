import { Gallery } from "react-grid-gallery";


export default function Filmstrip(props) {

    return (
        <Gallery 
            images={props.photos}
            enableImageSelection={props.enableSelection}
            onClick={props.onClick}
        />
    );
}