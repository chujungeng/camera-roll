import { Navigate,Outlet } from 'react-router-dom';
import Images from '../routes/Images';
import Image from '../routes/Image';
import Albums from '../routes/Albums';
import Album from '../routes/Album';
import Tags from '../routes/Tags';
import Tag from '../routes/Tag';
import Login from '../routes/Login';

function Layout() {

    return (
        <div>
            <Outlet />
        </div>
    );
}

function NotFound() {
    return (
        <main style={{ padding: "1rem", minHeight: "50vh" }}>
            <h2 className="text-light">404</h2>
            <p className="text-light">There's nothing here!</p>
        </main>
    );
}

const routes = (isLoggedIn) => [
    {
        path: "/",
        element: isLoggedIn? <Layout />: <Navigate to="/login" />,
        children: [
            { index: true, element: <Images /> },
            { path: "/images", element: <Images /> },
            { path: "/images/:imageID", element: <Image /> },
            { path: "/albums", element: <Albums /> },
            { path: "/albums/:albumID", element: <Album /> },
            { path: "/tags", element: <Tags /> },
            { path: "/tags/:tagID", element: <Tag /> },
            { path: "*", element: <NotFound /> }
        ],
    },
    { path: "/login", element: <Login /> },
];
  
export default routes;