'use client';

import { motion } from 'framer-motion';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { fadeIn, stagger } from '@/lib/animations';

export function FeaturesSection() {
  return (
    <section className='py-20 bg-light'>
      <motion.div 
        initial='initial'
        whileInView='animate'
        viewport={{ once: true }}
        variants={stagger}
        className='container max-w-6xl mx-auto px-4'
      >
        <div className='mx-auto mb-16 max-w-3xl text-center'>
          <motion.h2 
            variants={fadeIn}
            className='text-3xl font-bold tracking-tight text-primary sm:text-4xl lg:text-5xl'
          >
            Enterprise-Grade Risk Management
          </motion.h2>
          <motion.p 
            variants={fadeIn}
            className='mt-6 text-lg text-muted max-w-2xl mx-auto'
          >
            Comprehensive suite of tools powered by advanced machine learning algorithms
          </motion.p>
        </div>
        <div className='mx-auto grid max-w-5xl gap-8 md:grid-cols-3'>
          {[
            {
              title: 'Real-Time Analysis',
              description: 'Continuous monitoring and assessment',
              features: ['Live risk scoring', 'Market condition tracking', 'Automated alerts']
            },
            {
              title: 'Scenario Testing',
              description: 'Advanced market simulations',
              features: ['Stress testing', 'Bull market scenarios', 'Risk forecasting']
            },
            {
              title: 'ML Intelligence',
              description: 'Smart predictive analytics',
              features: ['Pattern recognition', 'Anomaly detection', 'Trend analysis']
            }
          ].map((feature, index) => (
            <motion.div key={index} variants={fadeIn}>
              <Card className='bg-card border-light shadow-card rounded-lg p-6'>
                <CardHeader>
                  <CardTitle className='card-title'>{feature.title}</CardTitle>
                  <CardDescription className='card-description'>
                    {feature.description}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <ul className='list-inside list-disc space-y-2 text-muted'>
                    {feature.features.map((item, i) => (
                      <li key={i}>{item}</li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </section>
  );
}
