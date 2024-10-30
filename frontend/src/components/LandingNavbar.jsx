import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider } from "@blueprintjs/core";
import { useNavigate } from "react-router-dom";

export const LandingNavbar = () => {
    const navigate = useNavigate();
    const handleClick = () => { navigate(`/account`); };

    return (
        <Navbar fixedToTop={true}>
            <Navbar.Group>
                <Icon icon="graph" size={25} />
                <Navbar.Divider />
                <Navbar.Heading><h2>College User Testing</h2></Navbar.Heading>
            </Navbar.Group>
            <Navbar.Group align={Alignment.RIGHT}>
                <Navbar.Divider />
                <Button className="bp5-minimal" icon="user" text="Account" onClick={handleClick} />
            </Navbar.Group>
        </Navbar>
    );
}