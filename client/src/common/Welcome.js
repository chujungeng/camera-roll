export default function Welcome(props) {
    return (
        <main className="pt-5" style={{flex: 1, minHeight: "55vh"}}>
            <h2 className="text-light">Welcome{props.user === null? null: `, ${props.user}`}</h2>
        </main>
    );
}