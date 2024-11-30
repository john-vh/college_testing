import { Button, Card, OverlayToaster, Classes, Checkbox, H2, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider } from "@blueprintjs/core";
import { useEffect } from "react";
import { PostingInfo } from "../hooks/useAllPostings";
import React from "react";
import { useApplyPosting } from "../hooks/useApplyPosting.ts";

interface PostingContentProps {
    post: PostingInfo
}

export const PostingContent = ({ post }: PostingContentProps) => {
    const myToaster = OverlayToaster.createAsync({ position: "bottom-right" });
    const applyPosting = useApplyPosting();

    const handleClick = (post) => {
        applyPosting(post, "");
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
                        <H2>{post.title}</H2>
                    </div>

                    <Button style={{ color: "black" }} disabled={true} minimal={true} outlined={true}>Virtual Live</Button>

                </div>
                <p><strong>Description</strong></p>
                <p>{post.desc}</p>
                <p><strong>Compensation Information</strong></p>
                <p><strong>${post.pay}</strong> via Paypal upon reviewed feedback completion, guaranteed within 7 business days</p>
                <div className='Footer'>
                    <div className="Flex">
                        <div className='icon-p'>
                            <Icon icon="user" size={30} />
                        </div>
                        <div>
                            <strong>Founder Name</strong>
                            <div>
                                Founder Description
                            </div>
                        </div>
                    </div>
                    <div style={{ padding: '10px', minWidth: '130px' }}>
                        <Button intent="primary" onClick={() => handleClick(post)}>I'm interested!</Button>
                    </div>
                </div>
            </Card>
        </div>
    );
}