import React, { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { LandingNavbar } from '../components/LandingNavbar.tsx';
import { PostingContent } from '../components/PostingContent.tsx';
import { useLocation } from 'react-router-dom';
import { PostingInfo } from '../hooks/useAllPostings.ts';
import { BusinessInfo } from '../hooks/useBusinessInfo.ts';

export const Posting = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const { post, businessMap } = location.state as { post: PostingInfo, businessMap: Map<string, BusinessInfo> } || {};

    useEffect(() => {
        if (!post) {
            navigate('/');
        }
    }, [post, navigate]);

    if (!post) { return; }

    return (
        <div>
            <LandingNavbar />
            <PostingContent post={post} businessMap={businessMap} />
        </div>
    );
};

export default Posting;
