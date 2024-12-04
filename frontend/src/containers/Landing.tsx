import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, InputGroup, NonIdealState } from "@blueprintjs/core";
import React, { useState, useMemo } from 'react';
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useNavigate } from 'react-router-dom';
import useAllPostings, { PostingInfo } from '../hooks/useAllPostings.ts';
import useAccountInfo from "../hooks/useAccountInfo.ts";

export const Landing = () => {
  const account = useAccountInfo();

  const handleLogin = () => {
    window.location.href = `${process.env.REACT_APP_API_URL}/auth/google`;
  }

  if (account == null) {
    return (
      <div>
        <LandingNavbar showHome={false} />
        <div className='App' style={{ margin: "50px" }}>
          <NonIdealState
            icon="log-in"
            title="Please log in"
            description="In order to use TestHive, please log in with a Google account"
            action={<Button onClick={() => handleLogin()}>Log in</Button>}
          />
        </div>
      </div>
    );
  }
  return (
    <div>
      <LandingNavbar />
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
    navigate(`/posting/${post.id}`, { state: { post } });
  };
  const [searchQuery, setSearchQuery] = useState("");

    const filteredPosts = useMemo(() => {
        return postingInfo?.filter(post =>
            post.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
            post.desc.toLowerCase().includes(searchQuery.toLowerCase()));
    }, [postingInfo, searchQuery]);

    if (postingInfo != null) {
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
                                    <p>{businessMap.get(post.business_id)?.name ?? "Startup Name"}</p>
                                </div>
                            </div>
                            <p className="bp5-text-muted">{post.desc}</p>
                            <div className='Flex align-right'>
                                <div className="gap">Compensation: ${post.pay}</div>
                                <Button intent="primary" onClick={() => handleClick(post)}>Details</Button>
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
