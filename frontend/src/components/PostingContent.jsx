import { Button, Card, OverlayToaster, Classes, Checkbox, H2, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider } from "@blueprintjs/core";
import { usePostingData } from "../hooks/usePostingData";
import { useEffect } from "react";

export const PostingContent = ({ id }) => {
    const data = usePostingData(id);

    const myToaster = OverlayToaster.createAsync({ position: "bottom-right" });

    const handleClick = () => {
        myToaster.then(toaster => toaster.show({ message: "Startup notified of interest!", intent: "success" }));
    }


    return (
        <div className="Posting">
            <Card interactive={false} >
                <div className="Flex" style={{ justifyContent: "space-between" }}>
                    <div className='Flex'>
                        <div className='icon-p'>
                            <Icon icon="bookmark" size={70} />
                        </div>
                        <H2>{data.name}</H2>
                    </div>

                    <Button style={{ color: "black" }} disabled={true} minimal={true} outlined={true}>Virtual Live</Button>

                </div>
                <p><strong>Description</strong></p>
                <p>{data.testDescription}</p>
                <p><strong>Testing Instructions</strong></p>
                <p>{data.instructions}</p>
                <p><strong>Compensation Information</strong></p>
                <p><strong>$5</strong> via Paypal upon reviewed feedback completion, guaranteed within 7 business days</p>
                <div className='Footer'>
                    <div className="Flex">
                        <div className='icon-p'>
                            <Icon icon="user" size={30} />
                        </div>
                        <div>
                            <strong>{data.founderName}</strong>
                            <div>
                                {data.startupDescription}
                            </div>
                        </div>
                    </div>
                    <div style={{ padding: '10px', minWidth: '130px' }}>
                        <Button intent="primary" onClick={handleClick}>I'm interested!</Button>
                    </div>
                </div>
            </Card>
        </div>
    );
}