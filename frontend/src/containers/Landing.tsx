import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, InputGroup } from "@blueprintjs/core";
import React, { useState, useMemo } from 'react';
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useNavigate } from 'react-router-dom';
import useAllPostings from '../hooks/useAllPostings.ts';

export const Landing = () => {
    return (
        <div>
            <LandingNavbar />
            <div className='App'>
                <FilterBar />
                <TestList />
            </div>
        </div>

    );
}

const TestList = () => {
    const posts = useAllPostings();

    const navigate = useNavigate();
    const handleClick = (id) => { navigate(`/posting/${id}`); };
    const [searchQuery, setSearchQuery] = useState("");

    const filteredPosts = useMemo(() => {
        return posts?.filter(post =>
            post.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
            post.desc.toLowerCase().includes(searchQuery.toLowerCase()));
    }, [posts, searchQuery]);

    if (posts != null) {
        return (
            <div className='Test-list'>
                <div style={{ position: "absolute" }}>
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
                                    <p>Startup Name</p>
                                </div>
                            </div>
                            <p className="bp5-text-muted">{post.desc}</p>
                            <div className='Flex align-right'>
                                <div className="gap">Compensation: $5</div>
                                <Button intent="primary" onClick={() => handleClick(post.id)}>Details</Button>
                            </div>
                        </Card>
                    </div>

                ))}
            </div>
        );
    }
}

const FilterBar = () => {
    const [isOpen, setOpen] = useState(true);

    if (isOpen) {
        return (
            <div className='Filter-bar'>
                <div className="Filter-header">
                    <div><strong>Filter Bar</strong></div>
                    <Button intent="primary" onClick={() => setOpen(!isOpen)} icon='filter'></Button>
                </div>
                <Checkbox label='Filter 1' />
                <Checkbox label='Filter 2' />
                <Checkbox label='Filter 3' />
            </div>
        );
    }

    else {
        return (
            <div className='hover'>
                <Button intent="primary" onClick={() => setOpen(!isOpen)} icon='filter'></Button>
            </div>

        )
    }
}