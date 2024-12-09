import { Button, Card, OverlayToaster, Classes, Checkbox, H2, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider, Tag, Colors } from "@blueprintjs/core";
import { useEffect } from "react";
import { PostingInfo } from "../hooks/useAllPostings";
import React from "react";
import { useApplyPosting } from "../hooks/useApplyPosting.ts";
import { formatDate } from "./UserApplicationInfo.tsx";
import { BusinessInfo } from "../hooks/useBusinessInfo.ts";
import { BackButton } from "./BackButton.tsx";

interface PostingContentProps {
    post: PostingInfo
    businessMap: Map<string, BusinessInfo>
}

export const PostingContent = ({ post, businessMap }: PostingContentProps) => {
    // const myToaster = OverlayToaster.createAsync({ position: "bottom-right" });
    const applyPosting = useApplyPosting();

    const handleClick = (post) => {
        applyPosting(post, "");
        // myToaster.then(toaster => toaster.show({ message: "Startup notified of interest!", intent: "success" }));
    }

    const businessInfo = businessMap.get(post.business_id);

    return (
        <div>
            <BackButton message={"Back to postings"} />
            <div className="Posting">
                <Card interactive={false} >
                    <div className="Flex" style={{ justifyContent: "space-between" }}>
                        <div className='Flex'>
                            <div className='icon-p'>
                                <Icon icon="bookmark" size={70} />
                            </div>
                            <div style={{ paddingLeft: "10px" }}>
                                <H2>{post.title}</H2>
                                <Tag minimal>{formatDate(post.created_at)}</Tag>
                            </div>
                        </div>
                    </div>
                    <p><strong>Description</strong></p>
                    <p>{post.desc}</p>
                    <p><strong>Compensation Information</strong></p>
                    <p><strong>${post.pay}</strong> via Paypal upon reviewed feedback completion, guaranteed within 7 business days</p>
                    <div className='Footer'>
                        <div className='Flex'>
                            <div className='icon-p'>
                                <Icon icon="user" size={35} />
                            </div>
                            <div style={{ display: "flex", flexDirection: "column", gap: "5px" }}>
                                <a href={businessInfo?.website}><strong>{businessInfo?.name}</strong></a>
                                <div>{businessInfo?.desc}</div>
                            </div>

                        </div>


                        <div style={{ padding: '10px', minWidth: '130px' }}>
                            <Button style={{ background: Colors.VIOLET2, color: Colors.WHITE }} onClick={() => handleClick(post)}>I'm interested!</Button>
                        </div>
                    </div>
                </Card>
            </div>
        </div>
    );
}