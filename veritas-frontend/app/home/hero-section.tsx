'use client';

import Link from 'next/link';
import { motion } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { fadeIn, stagger } from '@/lib/animations';

export function HeroSection() {
  return (
    <section className='relative min-h-[80vh] flex items-center justify-center py-20'>
      <div className='absolute inset-0 bg-[url("/grid.svg")] opacity-10 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]' />
      <motion.div 
        initial='initial'
        animate='animate'
        variants={stagger}
        className='container max-w-6xl mx-auto px-4'
      >
        <div className='flex flex-col items-center gap-8 text-center mx-auto max-w-3xl'>
          <motion.h1 
            variants={fadeIn}
            className='text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl lg:text-7xl text-white'
          >
            Real-Time Risk Assessment for{' '}
            <span className='bg-gradient-to-r from-white to-[#373b3b] bg-clip-text text-transparent'>
              Veritas Vault
            </span>
          </motion.h1>
          <motion.p 
            variants={fadeIn}
            className='mx-auto max-w-[700px] text-lg text-[#373b3b] md:text-xl'
          >
            Advanced ML-powered risk analysis and monitoring system for DeFi operations.
          </motion.p>
          <div className='mt-6'>
            <motion.div variants={fadeIn}>
              <Button 
                asChild 
                size='lg' 
                className='bg-[#373b3b] text-white hover:bg-black/80 transition-colors duration-300'
              >
                <Link href='/dashboard'>Launch Dashboard</Link>
              </Button>
            </motion.div>
          </div>
        </div>
      </motion.div>
    </section>
  );
}
