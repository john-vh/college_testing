import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, InputGroup, NonIdealState, Tag, Dialog, DialogBody, DialogFooter, Colors } from "@blueprintjs/core";
import React, { useState, useMemo } from 'react';
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useNavigate } from 'react-router-dom';
import useAllPostings, { PostingInfo } from '../hooks/useAllPostings.ts';
import useAccountInfo from "../hooks/useAccountInfo.ts";
import { formatDate } from "../components/UserApplicationInfo.tsx";
import { HelpDialog } from "../components/HelpDialog.tsx";
import { LoginSplash } from "../components/LoginSplash.tsx";

export const Landing = () => {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const account = useAccountInfo();

    const handleDialog = (isOpen: boolean) => {
        setIsOpen(isOpen);
    }

    const handleLogin = () => {
        window.location.href = `${process.env.REACT_APP_API_URL}/auth/google`;
    }

    if (account == null) {
        return (
            <LoginSplash handleDialog={handleDialog} handleLogin={handleLogin} />
        );
    }

    return (
        <div>
            <HelpDialog isOpen={isOpen} setIsOpen={setIsOpen} />
            <LandingNavbar handleDialog={handleDialog} />
            <div className='App'>
                <TestList />
            </div>
        </div>

    );
}

const TestList = () => {
    const { postingInfo, businessMap } = useAllPostings();
    const navigate = useNavigate();

    const handleClick = (post: PostingInfo) => {
        navigate(`/posting/${post.id}`, { state: { post, businessMap } });
    };
    const [searchQuery, setSearchQuery] = useState("");

    const filteredPosts = useMemo(() => {
        return postingInfo?.filter(post =>
            post.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
            post.desc.toLowerCase().includes(searchQuery.toLowerCase())).sort((a, b) => -a.created_at.localeCompare(b.created_at));
    }, [postingInfo, searchQuery]);

    if (postingInfo != null) {
        return (
            <div className='Test-list'>
                <div style={{ position: "absolute", width: "600px" }}>
                    <InputGroup
                        placeholder="Search..."
                        type="search"
                        value={searchQuery}
                        onValueChange={(value) => setSearchQuery(value)}
                    />
                </div>
                <div style={{ paddingBottom: "50px" }} />
                {filteredPosts != null && filteredPosts.map((post, index) => (
                    <div className='Card'>
                        <Card interactive={true} >
                            <div className='Flex'>
                                <div className='icon-p'>
                                    <Icon icon="bookmark" size={35}></Icon>
                                </div>
                                <div>
                                    <H5>{post.title}</H5>
                                    <p>{businessMap.get(post.business_id)?.name ?? "Startup Name"}</p>
                                </div>
                            </div>
                            <p className="bp5-text-muted">{post.desc}</p>
                            <div className="Flex" style={{ justifyContent: "space-between" }}>
                                <Tag minimal>{formatDate(post.created_at)}</Tag>
                                <div className='Flex'>
                                    <div className="gap">Compensation: ${post.pay}</div>
                                    <Button style={{ background: Colors.VIOLET2, color: Colors.WHITE }} onClick={() => handleClick(post)}>Details</Button>
                                </div>
                            </div>
                        </Card>
                    </div>

                ))}
            </div>
        );
    }
}